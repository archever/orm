package orm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ScanRow retore the sql result
type ScanRow struct {
	Value  string          // column value
	Column *sql.ColumnType // column info
	Valid  bool            // whether NULL
}

type M map[string]any

// Scan fetch the sql data
func Scan(rows *sql.Rows) ([]map[string]*ScanRow, error) {
	defer rows.Close()

	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	rets := []map[string]*ScanRow{}

	for rows.Next() {
		ret := map[string]*ScanRow{}
		scanI := make([]interface{}, len(cols))
		scanV := make([][]byte, len(cols))
		for i := range scanV {
			scanI[i] = &scanV[i]
		}
		err = rows.Scan(scanI...)
		if err != nil {
			return nil, err
		}

		for i := range cols {
			v := ""
			isValid := true
			if scanV[i] == nil {
				isValid = false
			} else {
				v = string(scanV[i])
			}
			ret[cols[i].Name()] = &ScanRow{
				Value:  v,
				Column: cols[i],
				Valid:  isValid,
			}
		}
		rets = append(rets, ret)
	}
	return rets, nil
}

// ToValue return the default value for database type
func (r *ScanRow) ToValue() (interface{}, error) {
	var ret interface{}
	var err error
	switch r.Column.DatabaseTypeName() {
	case "INT", "SMALLINT", "TINYINT":
		ret, err = r.ToInt()
	case "BIGINT":
		ret, err = r.ToInt64()
	case "FLOAT", "DECIMAL":
		ret, err = r.ToFloat64()
	case "DATETIME", "TIMESTAMP", "TIME", "DATE":
		ret, err = r.ToTime()
	case "VARCHAR", "CHAR", "TEXT":
		ret, err = r.ToString()
	case "BLOB":
		ret, err = r.ToBytes()
	default:
		Debug.Printf("unknown column: %s", r.Column.DatabaseTypeName())
		ret, err = r.ToString()
	}
	return ret, err
}

// String string
func (r *ScanRow) String() string {
	return r.Value
}

// ToString string
func (r *ScanRow) ToString() (string, error) {
	if !r.Valid {
		return "", ErrNull
	}
	return r.Value, nil
}

// ToBytes []byte
func (r *ScanRow) ToBytes() ([]byte, error) {
	if !r.Valid {
		return nil, ErrNull
	}
	return []byte(r.Value), nil
}

// ToInt64 int64
func (r *ScanRow) ToInt64() (int64, error) {
	if !r.Valid {
		return 0, ErrNull
	}
	v, err := strconv.ParseInt(r.Value, 10, 64)
	return v, err
}

// ToFloat64 float64
func (r *ScanRow) ToFloat64() (float64, error) {
	if !r.Valid {
		return 0, ErrNull
	}
	v, err := strconv.ParseFloat(r.Value, 64)
	return v, err
}

// ToFloat32 float32
func (r *ScanRow) ToFloat32() (float32, error) {
	if !r.Valid {
		return 0, ErrNull
	}
	v, err := strconv.ParseFloat(r.Value, 32)
	return float32(v), err
}

// ToInt int
func (r *ScanRow) ToInt() (int, error) {
	if !r.Valid {
		return 0, ErrNull
	}
	v, err := strconv.Atoi(r.Value)
	return v, err
}

// ToInt8 int8
func (r *ScanRow) ToInt8() (int8, error) {
	if !r.Valid {
		return 0, ErrNull
	}
	v, err := strconv.Atoi(r.Value)
	return int8(v), err
}

// ToByte byte
func (r *ScanRow) ToByte() (byte, error) {
	if !r.Valid {
		return 0, ErrNull
	}
	v, err := strconv.Atoi(r.Value)
	return byte(v), err
}

// ToBool bool
func (r *ScanRow) ToBool() (bool, error) {
	if !r.Valid {
		return false, ErrNull
	}
	switch strings.ToLower(r.Value) {
	case "", "0", "false":
		return false, nil
	default:
		return true, nil
	}
}

// ToTime *time.Time
func (r *ScanRow) ToTime() (*time.Time, error) {
	if !r.Valid {
		return nil, ErrNull
	}
	v, err := time.ParseInLocation("2006-01-02 15:04:05", r.Value, time.Now().Location())
	if err == nil {
		return &v, nil
	}
	v, err = time.ParseInLocation("2006-01-02", r.Value, time.Now().Location())
	if err == nil {
		return &v, nil
	}
	return nil, err
}

// ToJSON json
func (r *ScanRow) ToJSON(dest interface{}) error {
	if !r.Valid {
		return ErrNull
	}
	body, err := r.ToBytes()
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, dest)
	return err
}

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

func ScanQueryFields(dest []FieldIfc, rows *sql.Rows) error {
	fields, err := rows.Columns()
	if err != nil {
		return err
	}
	fieldMap := map[string]int{}
	for i, field := range fields {
		// 名称一样的场景, 跨表查询
		if _, ok := fieldMap[field]; ok {
			return ScanQueryFieldsWithOrder(dest, rows)
		}
		fieldMap[field] = i
	}
	values := make([]interface{}, len(fields))
	for _, field := range dest {
		index, ok := fieldMap[field.ColName(false)]
		if !ok {
			return fmt.Errorf("field not found: %s", field.ColName(false))
		}
		values[index] = field.RefVal()
	}
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}
	}
	for _, field := range dest {
		field.setPreVal(field.Val())
	}
	return nil
}

func ScanQueryFieldsWithOrder(dest []FieldIfc, rows *sql.Rows) error {
	values := make([]any, 0, len(dest))
	for _, field := range dest {
		values = append(values, field.RefVal())
	}
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}
	}
	for _, field := range dest {
		field.setPreVal(field.Val())
	}
	return nil
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
	if !ok {
		fieldName, ok = field.Tag.Lookup("json")
	}
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
		if _, ok = row[fieldName]; !ok {
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
