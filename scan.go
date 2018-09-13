// scan countians fetching sql data method and data type convert

package orm

import (
	"database/sql"
	"log"
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
	return rets, nil
}

// ToValue return the default value for database type
func (r *ScanRow) ToValue() (interface{}, error) {
	log.Printf("name: %-10s type: %-10s value: %s", r.Column.Name(), r.Column.DatabaseTypeName(), string(r.Value))
	var ret interface{}
	var err error
	switch r.Column.DatabaseTypeName() {
	case "INT":
		ret, err = r.ToInt64()
	case "DATETIME", "DATE", "TIMESTAMP":
		ret, err = r.ToTime()
	case "VARCHAR":
		ret = r.ToString()
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
	v, err := strconv.ParseInt(r.ToString(), 10, 64)
	return v, err
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
