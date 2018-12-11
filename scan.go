package orm

import (
	"database/sql"
	"encoding/json"
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
		rows.Scan(scanI...)

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
