package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Decimal2 int

func (n *Decimal2) Scan(value interface{}) (err error) {

	if value == nil {
		return
	}

	var valStr string
	switch value.(type) {
	case []byte:
		valStr = string(value.([]byte))
	case string:
		valStr = value.(string)
	default:
		err = ErrInvalidType
		return
	}

	var val float64
	if val, err = strconv.ParseFloat(valStr, 64); err != nil {
		return
	}

	*n = Decimal2(val * 100)

	return
}

func (n Decimal2) Value() (driver.Value, error) {

	res := fmt.Sprintf("%d.%02d", n/100, n%100)

	return res, nil

}
