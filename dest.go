package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Unmarshaler data to read in sql field
type Unmarshaler interface {
	UnmarshalSQL(*ScanRow) error
}

// Marshaler data to store in sql field
type Marshaler interface {
	MarshalSQL() (string, error)
}

// ITOMarshaler 判断 是否实现 Marshaler 接口
func ITOMarshaler(v reflect.Value) (Marshaler, bool) {
	var m Marshaler
	var ok bool
	if v.CanAddr() {
		m, ok = v.Addr().Interface().(Marshaler)
	}
	if !ok {
		m, ok = v.Interface().(Marshaler)
	}
	vl := v
	for !ok {
		if vl.Type().Kind() == reflect.Ptr && !vl.IsNil() {
			vl = vl.Elem()
			m, ok = vl.Interface().(Marshaler)
		} else {
			ok = false
			break
		}
	}
	return m, ok
}

// ITOUnmarshaler 判断 是否实现 Unmarshaler 接口
func ITOUnmarshaler(v reflect.Value) (Unmarshaler, bool) {
	var m Unmarshaler
	var ok bool
	if v.CanAddr() {
		m, ok = v.Addr().Interface().(Unmarshaler)
	}
	if !ok {
		m, ok = v.Interface().(Unmarshaler)
	}
	vl := v
	for !ok {
		if vl.Type().Kind() == reflect.Ptr && !vl.IsNil() {
			vl = vl.Elem()
			m, ok = vl.Interface().(Unmarshaler)
		} else {
			ok = false
			break
		}
	}
	return m, ok
}

func ScanQueryRows(dest interface{}, rows *sql.Rows) error {
	var err error
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("dest type not match: %T", dest)
	}
	// fetch data
	reses, err := Scan(rows)
	if err != nil {
		return err
	}
	if len(reses) == 0 {
		return nil
	}
	rv = indirect(rv)
	switch rv.Kind() {
	case reflect.Interface:
		if rv.NumMethod() == 0 {
			rv = reflect.ValueOf([]M{})
			rv, err := ToSlice(reses, rv)
			if err != nil {
				return err
			}
			reflect.ValueOf(dest).Elem().Set(rv)
		} else {
			return fmt.Errorf("dest type not match")
		}
	case reflect.Array:
		rv, err = ToArray(reses, rv)
		if err != nil {
			return err
		}
		reflect.ValueOf(dest).Elem().Set(rv)
	case reflect.Slice:
		rv, err = ToSlice(reses, rv)
		if err != nil {
			return err
		}
		reflect.ValueOf(dest).Elem().Set(rv)
	default:
		return fmt.Errorf("dest type not match: %T", dest)
	}
	return nil
}

func ScanQueryOne(dest interface{}, rows *sql.Rows) error {
	var err error
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("dest type not match: %T", dest)
	}
	// fetch one data
	reses, err := Scan(rows)
	if err != nil {
		return err
	}
	if len(reses) == 0 {
		return nil
	}
	res := reses[0]

	rv = indirect(rv)
	switch rv.Kind() {
	case reflect.Interface:
		if rv.NumMethod() == 0 {
			rv = reflect.ValueOf(M{})
			err := ToMap(res, rv)
			if err != nil {
				return err
			}
			reflect.ValueOf(dest).Elem().Set(rv)
			return nil
		}
	case reflect.Map:
		err = ToMap(res, rv)
	case reflect.Struct:
		err = ToStruct(res, rv)
	default:
		return fmt.Errorf("dest type not match: %T", dest)
	}
	return err
}

func GetFieldName(field reflect.StructField) (string, bool, bool) {
	fieldName, ok := field.Tag.Lookup("column")
	isOmitempty := false
	isIgnore := false
	// 判断是否是导出字段
	first := string(field.Name[0])
	if first == strings.ToLower(first) {
		isIgnore = true
	} else if ok {
		fieldNames := strings.Split(fieldName, ",")
		if len(fieldNames) == 1 {
			if fieldNames[0] == "omitempty" {
				isOmitempty = true
				fieldName = ""
			} else if fieldNames[0] == "-" {
				isIgnore = true
				fieldName = ""
			} else {
				fieldName = fieldNames[0]
			}
		} else {
			for _, f := range fieldNames {
				if f == "omitempty" {
					isOmitempty = true
					continue
				}
				if f == "-" {
					isIgnore = true
					fieldName = ""
					break
				}
				fieldName = f
			}
		}
	}
	if fieldName == "" {
		fieldName = field.Name
	}
	return fieldName, isOmitempty, isIgnore
}

// TODO 重新实现
func indirect(v reflect.Value) reflect.Value {
	v0 := v
	haveAddr := false
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && e.Elem().Kind() == reflect.Ptr {
				haveAddr = false
				v = e
				continue
			}
		}
		if v.Kind() != reflect.Ptr {
			break
		}
		if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Ptr {
			v = v.Elem()
			continue
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Elem().Kind() != reflect.Ptr && v.CanSet() {
			v = v.Elem()
			continue
		}
		if haveAddr {
			v = v0
			haveAddr = false
		} else if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}
	return v
}

// ToSlice rows to slice
func ToSlice(reses []map[string]*ScanRow, rv reflect.Value) (reflect.Value, error) {
	rt := rv.Type().Elem()
	isPtr := false
	if rt.Kind() == reflect.Ptr {
		rt = rv.Type().Elem().Elem()
		isPtr = true
	}
	switch rt.Kind() {
	case reflect.Map:
		for _, res := range reses {
			v := reflect.MakeMap(rt)
			err := ToMap(res, v)
			if err != nil {
				return rv, err
			}
			if isPtr {
				rv = reflect.Append(rv, v.Addr())
			} else {
				rv = reflect.Append(rv, v)
			}
		}
	case reflect.Struct:
		for _, res := range reses {
			v := reflect.New(rt)
			v = v.Elem()
			err := ToStruct(res, v)
			if err != nil {
				return rv, err
			}
			if isPtr {
				rv = reflect.Append(rv, v.Addr())
			} else {
				rv = reflect.Append(rv, v)
			}
		}
	default:
		return rv, fmt.Errorf("dest type not match")
	}
	return rv, nil
}

// ToArray rows to slice
func ToArray(reses []map[string]*ScanRow, rv reflect.Value) (reflect.Value, error) {
	rt := rv.Type().Elem()
	isPtr := false
	if rt.Kind() == reflect.Ptr {
		rt = rv.Type().Elem().Elem()
		isPtr = true
	}
	length := rv.Len()
	if length == 0 {
		return rv, nil
	}
	switch rt.Kind() {
	case reflect.Map:
		index := 0
		for _, res := range reses {
			if index >= length {
				break
			}
			v := reflect.MakeMap(rt)
			err := ToMap(res, v)
			if err != nil {
				return rv, err
			}
			if isPtr {
				rv.Index(index).Set(v.Addr())
			} else {
				rv.Index(index).Set(v)
			}
			index++
		}
	case reflect.Struct:
		index := 0
		for _, res := range reses {
			if index >= length {
				break
			}
			v := reflect.New(rt)
			v = v.Elem()
			err := ToStruct(res, v)
			if err != nil {
				return rv, err
			}
			if isPtr {
				rv.Index(index).Set(v.Addr())
			} else {
				rv.Index(index).Set(v)
			}
			index++
		}
	default:
		return rv, fmt.Errorf("dest type not match")
	}
	return rv, nil
}

// ToMap rows to map
func ToMap(row map[string]*ScanRow, rv reflect.Value) error {
	// assert key string type
	r0 := rv
	if r0.IsNil() {
		rv = reflect.MakeMap(rv.Type())
	}
	if rv.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("dest type not match")
	}
	for k, v := range row {
		var data interface{}
		var err error
		switch rv.Type().Elem().Kind() {
		case reflect.Interface:
			data, err = v.ToValue()
		case reflect.String:
			data, err = v.ToString()
		case reflect.Int64:
			data, err = v.ToInt64()
		case reflect.Uint8:
			data, err = v.ToByte()
		case reflect.Int:
			data, err = v.ToInt()
		case reflect.Int8:
			data, err = v.ToInt8()
		case reflect.Bool:
			data, err = v.ToBool()
		default:
			return fmt.Errorf("dest type not match")
		}
		if err == ErrNull {
			rv.SetMapIndex(reflect.ValueOf(k), reflect.Zero(rv.Type().Elem()))
			continue
		} else if err != nil {
			return err
		} else {
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(data))
		}
	}
	if r0.IsNil() {
		r0.Set(rv)
	}
	return nil
}

// ToStruct row to struct
func ToStruct(row map[string]*ScanRow, rv reflect.Value) error {
	for i := 0; i < rv.NumField(); i++ {
		ele := rv.Field(i)
		fieldName, _, isIgnore := GetFieldName(rv.Type().Field(i))
		if isIgnore {
			continue
		}
		var data *ScanRow
		var ok bool
		if data, ok = row[fieldName]; !ok {
			fieldName = strings.ToLower(fieldName)
		}
		if data, ok = row[fieldName]; !ok {
			continue
		}
		if ele.Type().Kind() == reflect.Ptr && ele.IsNil() {
			if !data.Valid {
				ele.Set(reflect.Zero(ele.Type()))
				continue
			}
			ele.Set(reflect.New(ele.Type().Elem()))
			ele = ele.Elem()
		}
		if m, ok := ITOUnmarshaler(ele); ok {
			err := m.UnmarshalSQL(data)
			if err == ErrNull {
				ele.Set(reflect.Zero(ele.Type()))
			} else if err != nil {
				return err
			} else {
				ele.Set(reflect.ValueOf(m).Elem())
			}
			continue
		}
		switch ele.Type().Kind() {
		case reflect.String:
			v := data.Value
			rv := reflect.ValueOf(v)
			if ele.Type().ConvertibleTo(rv.Type()) {
				rv = rv.Convert(ele.Type())
			}
			if !ele.Type().AssignableTo(rv.Type()) {
				return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
			}
			ele.Set(rv)
		case reflect.Float64:
			v, err := data.ToFloat64()
			if err != nil && err != ErrNull {
				return err
			}
			rv := reflect.ValueOf(v)
			if ele.Type().ConvertibleTo(rv.Type()) {
				rv = rv.Convert(ele.Type())
			}
			if !ele.Type().AssignableTo(rv.Type()) {
				return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
			}
			ele.Set(rv)
		case reflect.Int64:
			v, err := data.ToInt64()
			if err != nil && err != ErrNull {
				return err
			}
			rv := reflect.ValueOf(v)
			if ele.Type().ConvertibleTo(rv.Type()) {
				rv = rv.Convert(ele.Type())
			}
			if !ele.Type().AssignableTo(rv.Type()) {
				return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
			}
			ele.Set(rv)
		case reflect.Int:
			v, err := data.ToInt()
			if err != nil && err != ErrNull {
				return err
			}
			rv := reflect.ValueOf(v)
			if ele.Type().ConvertibleTo(rv.Type()) {
				rv = rv.Convert(ele.Type())
			}
			if !ele.Type().AssignableTo(rv.Type()) {
				return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
			}
			ele.Set(rv)
		case reflect.Int8:
			v, err := data.ToInt8()
			if err != nil && err != ErrNull {
				return err
			}
			rv := reflect.ValueOf(v)
			if ele.Type().ConvertibleTo(rv.Type()) {
				rv = rv.Convert(ele.Type())
			}
			if !ele.Type().AssignableTo(rv.Type()) {
				return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
			}
			ele.Set(rv)
		case reflect.Uint8:
			v, err := data.ToByte()
			if err != nil && err != ErrNull {
				return err
			}
			rv := reflect.ValueOf(v)
			if ele.Type().ConvertibleTo(rv.Type()) {
				rv = rv.Convert(ele.Type())
			}
			if !ele.Type().AssignableTo(rv.Type()) {
				return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
			}
			ele.Set(rv)
		case reflect.Bool:
			v, _ := data.ToBool()
			rv := reflect.ValueOf(v)
			if ele.Type().ConvertibleTo(rv.Type()) {
				rv = rv.Convert(ele.Type())
			}
			if !ele.Type().AssignableTo(rv.Type()) {
				return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
			}
			ele.Set(rv)
		case reflect.Ptr, reflect.Struct, reflect.Slice:
			if ele.Type().Kind() == reflect.Ptr {
				ele.Set(reflect.New(ele.Type().Elem()))
				ele = ele.Elem()
			}
			switch ele.Interface().(type) {
			case time.Time:
				v, err := data.ToTime()
				if err == ErrNull {
					ele.Set(reflect.Zero(ele.Type()))
					continue
				} else if err != nil {
					return err
				}
				ele.Set(reflect.ValueOf(*v))
			case []byte:
				v, err := data.ToBytes()
				if err == ErrNull {
					ele.Set(reflect.Zero(ele.Type()))
					continue
				} else if err != nil {
					return err
				}
				rv := reflect.ValueOf(v)
				if ele.Type().ConvertibleTo(rv.Type()) {
					rv = rv.Convert(ele.Type())
				}
				if !ele.Type().AssignableTo(rv.Type()) {
					return fmt.Errorf("%T not assignable to %T", ele.Interface(), rv.Interface())
				}
				ele.Set(rv)
			default:
				return fmt.Errorf("unknown type %T", ele.Interface())
			}
		default:
			return fmt.Errorf("unknown type %T", ele.Interface())
		}
	}
	return nil
}

// IToMap interface to map
func IToMap(v reflect.Value) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	// isPtr := false
	rv := v
	if v.Type().Kind() == reflect.Ptr {
		// isPtr = true
		rv = v.Elem()
	}
	switch rv.Type().Kind() {
	case reflect.Map:
		if rv.Len() == 0 {
			return nil, fmt.Errorf("save empty data")
		}
		keys := rv.MapKeys()
		for _, k := range keys {
			vi := rv.MapIndex(k)
			data := vi.Interface()
			if m, ok := ITOMarshaler(vi); ok {
				var err error
				data, err = m.MarshalSQL()
				if err == ErrNull {
					data = nil
				} else if err != nil {
					return ret, err
				}
			}
			ret[k.String()] = data
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Type().Field(i)
			fieldName, isOmitempty, isIgnore := GetFieldName(field)
			if isIgnore || (isOmitempty && IsEmptyValue(rv.Field(i))) {
				continue
			}
			vi := rv.Field(i)
			data := vi.Interface()
			if m, ok := ITOMarshaler(vi); ok {
				var err error
				data, err = m.MarshalSQL()
				if err == ErrNull {
					data = nil
				} else if err != nil {
					return ret, err
				}
			}
			ret[fieldName] = data
		}
	default:
		return ret, fmt.Errorf("invalid data type: %T", rv.Interface())
	}
	return ret, nil
}

func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
