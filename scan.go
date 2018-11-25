// scan countians fetching sql data method and data type convert

package orm

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

// ScanRow retore the sql result
type ScanRow struct {
	Value  []byte          // 数据值
	Column *sql.ColumnType // 数据名
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
			ret[cols[i].Name()] = &ScanRow{
				Value:  scanV[i],
				Column: cols[i],
			}
		}
		rets = append(rets, ret)
	}
	if len(rets) == 0 {
		return rets, ErrNotFund
	}
	return rets, nil
}

// ToValue return the default value for database type
func (r *ScanRow) ToValue() (interface{}, error) {
	var ret interface{}
	var err error
	switch r.Column.DatabaseTypeName() {
	case "INT", "SMALLINT", "TINYINT", "BIGINT":
		ret, err = r.ToInt64()
	case "DATETIME", "TIMESTAMP", "TIME":
		ret, err = r.ToTime()
	case "DATE":
		ret, err = r.ToDateTime()
	case "VARCHAR", "CHAR", "TEXT":
		ret = r.ToString()
	case "BLOB":
		ret = r.Value
	default:
		ret = r.ToString()
	}
	return ret, err
}

// ToString return the string value
func (r *ScanRow) ToString() string {
	return string(r.Value)
}

// ToInt64 return the int64 value
func (r *ScanRow) ToInt64() (int64, error) {
	if r.ToString() == "" {
		return 0, nil
	}
	v, err := strconv.ParseInt(r.ToString(), 10, 64)
	return v, err
}

// ToFloat64 return the float64 value
func (r *ScanRow) ToFloat64() (float64, error) {
	if r.ToString() == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(r.ToString(), 64)
	return v, err
}

// ToFloat32 return the float32 value
func (r *ScanRow) ToFloat32() (float32, error) {
	if r.ToString() == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(r.ToString(), 32)
	return float32(v), err
}

// ToInt return the int64 value
func (r *ScanRow) ToInt() (int, error) {
	if r.ToString() == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(r.ToString())
	return v, err
}

// ToInt8 return the int64 value
func (r *ScanRow) ToInt8() (int8, error) {
	if r.ToString() == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(r.ToString())
	return int8(v), err
}

// ToByte return the int64 value
func (r *ScanRow) ToByte() (byte, error) {
	if r.ToString() == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(r.ToString())
	return byte(v), err
}

// ToBool return the bool value
func (r *ScanRow) ToBool() bool {
	switch strings.ToLower(r.ToString()) {
	case "", "0", "false":
		return false
	default:
		return true
	}
}

// ToTime return the time value
func (r *ScanRow) ToTime() (time.Time, error) {
	var ret time.Time
	strV := r.ToString()
	if strV == "" {
		return ret, nil
	}
	v, err := time.ParseInLocation("2006-01-02 15:04:05", r.ToString(), time.Now().Location())
	return v, err
}

// ToDateTime return the time value
func (r *ScanRow) ToDateTime() (time.Time, error) {
	var ret time.Time
	strV := r.ToString()
	if strV == "" {
		return ret, nil
	}
	v, err := time.ParseInLocation("2006-01-02", r.ToString(), time.Now().Location())
	return v, err
}
