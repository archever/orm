// dest help to set data to destination pointer

package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// UnMarshaler data to read in sql field
type UnMarshaler interface {
	UnMarshalSQL([]byte) error
}

// Marshaler data to store in sql field
type Marshaler interface {
	MarshalSQL() (string, error)
}

func scanQueryRows(dest interface{}, rows *sql.Rows) error {
	var err error
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid dest type")
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
			return errors.New("should not be method interface")
		}
	case reflect.Array:
		// TODO: handler Array type
	case reflect.Slice:
		// ToSlice
		rv, err = ToSlice(reses, rv)
		if err != nil {
			return err
		}
		reflect.ValueOf(dest).Elem().Set(rv)
	default:
		return errors.New("invelid dest type")
	}
	return nil
}

func scanQueryOne(dest interface{}, rows *sql.Rows) error {
	var err error
	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("invalid dest type")
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
		// assert not interface method
		// set default value &map[string]interface{}{}
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
		// ToMap
		err = ToMap(res, rv)
	case reflect.Struct:
		// ToStruct
		err = ToStruct(res, rv)
	default:
		return errors.New("invelid dest type")
	}
	return err
}

func getFieldName(field reflect.StructField) string {
	fieldName, ok := field.Tag.Lookup("column")
	if ok {
		return fieldName
	}
	fieldName, ok = field.Tag.Lookup("json")
	if ok {
		fieldName = strings.Split(fieldName, ",")[0]
		return fieldName
	}
	fieldName = strings.ToLower(field.Name)
	return fieldName
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
		return rv, errors.New("invalid type while parsing to slice item")
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
		return errors.New("invalid map key, string expected")
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
	case reflect.Bool:
		for k, v := range row {
			rv.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v.ToBool()))
		}
	default:
		return errors.New("invalid dest map value type")
	}
	if r0.IsNil() {
		r0.Set(rv)
	}
	return nil
}

// ToStruct row to struct
func ToStruct(row map[string]*ScanRow, rv reflect.Value) error {
	for i := 0; i < rv.NumField(); i++ {
		// struct 值
		ele := rv.Field(i)
		// 字段名
		fieldName := getFieldName(rv.Type().Field(i))
		// 读取字段值
		if data, ok := row[fieldName]; ok {
			switch ele.Type().Kind() {
			case reflect.String:
				v := data.ToString()
				ele.Set(reflect.ValueOf(v))

			case reflect.Int64:
				v, err := data.ToInt64()
				if err != nil {
					return err
				}
				ele.Set(reflect.ValueOf(v))
			// check strcut type
			case reflect.Struct:
				// handler UnMarshalSQL
				if m, ok := ele.Addr().Interface().(UnMarshaler); ok {
					err := m.UnMarshalSQL(data.Value)
					if err != nil {
						return err
					}
					ele.Set(reflect.ValueOf(m).Elem())
					continue
				}
				// handler orther type
				switch ele.Interface().(type) {
				case time.Time:
					v, err := data.ToTime()
					if err != nil {
						return err
					}
					ele.Set(reflect.ValueOf(v))
				default:
					return fmt.Errorf("unknown type %T", ele.Interface())
				}
			}
		}
	}
	return nil
}

// IToMap interface to map
func IToMap(item interface{}) (map[string]interface{}, error) {
	v := reflect.ValueOf(item)
	ret := map[string]interface{}{}
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Type().Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for _, k := range keys {
			ret[k.String()] = v.MapIndex(k).Interface()
		}
	case reflect.Struct:
		// check seriale interface
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
			ret[getFieldName(field)] = data
		}
	default:
		return ret, errors.New("not a valid data")
	}
	return ret, nil
}
