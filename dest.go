// dest help to set data to destination pointer

package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// UnMarshaler data to read in sql field
type UnMarshaler interface {
	UnMarshalSQL(*ScanRow) error
}

// Marshaler data to store in sql field
type Marshaler interface {
	MarshalSQL() (string, error)
}

// ITOMarshaler TOOD: 这个方法不能判断指针接受者
func ITOMarshaler(m interface{}) Marshaler {
	var ret Marshaler
	if mi, ok := m.(Marshaler); ok {
		ret = mi
	}
	return ret
}

func ScanQueryRows(dest interface{}, rows *sql.Rows) error {
	var err error
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrDestType
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
		// assert not interface method
		// set default value &[]map[string]interface{}{}
		if rv.NumMethod() == 0 {
			rv = reflect.ValueOf([]map[string]interface{}{})
			rv, err := ToSlice(reses, rv)
			if err != nil {
				return err
			}
			reflect.ValueOf(dest).Elem().Set(rv)
		} else {
			return ErrDestType
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
		return ErrDestType
	}
	return nil
}

func ScanQueryOne(dest interface{}, rows *sql.Rows) error {
	var err error
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrDestType
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
			rv = reflect.ValueOf(map[string]interface{}{})
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
		return ErrDestType
	}
	return err
}

func getFieldName(field reflect.StructField) (string, bool) {
	fieldName, ok := field.Tag.Lookup("column")
	isOmitempty := false
	if ok {
		fieldNames := strings.Split(fieldName, ",")
		if len(fieldNames) == 1 {
			if fieldNames[0] == "omitempty" {
				isOmitempty = true
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
				fieldName = f
			}
		}
	}
	if fieldName == "" {
		fieldName = strings.ToLower(field.Name)
	}
	return fieldName, isOmitempty
}

func indirect(v reflect.Value) reflect.Value {
	v0 := v
	haveAddr := false
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	if v.Kind() == reflect.Interface && !v.IsNil() {
		e := v.Elem()
		if e.Kind() == reflect.Ptr && !e.IsNil() && e.Elem().Kind() == reflect.Ptr {
			haveAddr = false
			v = e
		}
	}
	if v.Kind() != reflect.Ptr {
		return v
	}

	if v.Elem().Kind() != reflect.Ptr && v.CanSet() {
		return v
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	if haveAddr {
		v = v0
		haveAddr = false
	} else {
		v = v.Elem()
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
		return rv, ErrDestType
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
		return rv, ErrDestType
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
		return ErrDestType
	}
	switch rv.Type().Elem().Kind() {
	case reflect.Interface:
		for k, v := range row {
			data, err := v.ToValue()
			if err != nil {
				return err
			}
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(data))
		}
	case reflect.String:
		for k, v := range row {
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v.ToString()))
		}
	case reflect.Int64:
		for k, v := range row {
			data, err := v.ToInt64()
			if err != nil {
				return err
			}
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(data))
		}
	case reflect.Uint8:
		for k, v := range row {
			data, err := v.ToByte()
			if err != nil {
				return err
			}
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(data))
		}
	case reflect.Int:
		for k, v := range row {
			data, err := v.ToInt()
			if err != nil {
				return err
			}
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(data))
		}
	case reflect.Int8:
		for k, v := range row {
			data, err := v.ToInt8()
			if err != nil {
				return err
			}
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(data))
		}
	case reflect.Bool:
		for k, v := range row {
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v.ToBool()))
		}
	default:
		return ErrDestType
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
		fieldName, _ := getFieldName(rv.Type().Field(i))
		if data, ok := row[fieldName]; ok {

			if m, ok := ele.Addr().Interface().(UnMarshaler); ok {
				err := m.UnMarshalSQL(data)
				if err != nil {
					return err
				}
				ele.Set(reflect.ValueOf(m).Elem())
				continue
			}

			switch ele.Type().Kind() {
			case reflect.String:
				v := data.ToString()
				rv := reflect.ValueOf(v)
				if ele.Type().ConvertibleTo(rv.Type()) {
					rv = rv.Convert(ele.Type())
				}
				if !ele.Type().AssignableTo(rv.Type()) {
					return ErrNotAssignable
				}
				ele.Set(rv)
			case reflect.Int64:
				v, err := data.ToInt64()
				if err != nil {
					return err
				}
				rv := reflect.ValueOf(v)
				if ele.Type().ConvertibleTo(rv.Type()) {
					rv = rv.Convert(ele.Type())
				}
				if !ele.Type().AssignableTo(rv.Type()) {
					return ErrNotAssignable
				}
				ele.Set(rv)
			case reflect.Int:
				v, err := data.ToInt()
				if err != nil {
					return err
				}
				rv := reflect.ValueOf(v)
				if ele.Type().ConvertibleTo(rv.Type()) {
					rv = rv.Convert(ele.Type())
				}
				if !ele.Type().AssignableTo(rv.Type()) {
					return ErrNotAssignable
				}
				ele.Set(rv)
			case reflect.Int8:
				v, err := data.ToInt8()
				if err != nil {
					return err
				}
				rv := reflect.ValueOf(v)
				if ele.Type().ConvertibleTo(rv.Type()) {
					rv = rv.Convert(ele.Type())
				}
				if !ele.Type().AssignableTo(rv.Type()) {
					return ErrNotAssignable
				}
				ele.Set(rv)
			case reflect.Uint8:
				v, err := data.ToByte()
				if err != nil {
					return err
				}
				rv := reflect.ValueOf(v)
				if ele.Type().ConvertibleTo(rv.Type()) {
					rv = rv.Convert(ele.Type())
				}
				if !ele.Type().AssignableTo(rv.Type()) {
					return ErrNotAssignable
				}
				ele.Set(rv)
			case reflect.Bool:
				v := data.ToBool()
				rv := reflect.ValueOf(v)
				if ele.Type().ConvertibleTo(rv.Type()) {
					rv = rv.Convert(ele.Type())
				}
				if !ele.Type().AssignableTo(rv.Type()) {
					return ErrNotAssignable
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
					if err != nil {
						return err
					}
					ele.Set(reflect.ValueOf(v))
				case []byte:
					rv := reflect.ValueOf(data.Value)
					if ele.Type().ConvertibleTo(rv.Type()) {
						rv = rv.Convert(ele.Type())
					}
					if !ele.Type().AssignableTo(rv.Type()) {
						return ErrNotAssignable
					}
					ele.Set(rv)
				default:
					return fmt.Errorf("unknown type %T", ele.Interface())
				}
			default:
				return fmt.Errorf("unknown type %T", ele.Interface())
			}
		}
	}
	return nil
}

// IToMap interface to map
func IToMap(v reflect.Value) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Type().Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for _, k := range keys {
			data := v.MapIndex(k).Interface()
			if m, ok := data.(Marshaler); ok {
				var err error
				data, err = m.MarshalSQL()
				if err != nil {
					return ret, err
				}
			}
			ret[k.String()] = data
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			data := v.Field(i).Interface()
			if m, ok := data.(Marshaler); ok {
				var err error
				data, err = m.MarshalSQL()
				if err != nil {
					return ret, err
				}
			}
			fieldName, isOmitempty := getFieldName(field)
			if isOmitempty && isEmptyValue(v.Field(i)) {
				continue
			}
			ret[fieldName] = data
		}
	default:
		return ret, ErrDestType
	}
	return ret, nil
}

func isEmptyValue(v reflect.Value) bool {
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
