package types

import (
	"testing"
	"fmt"
)

func TestDecimal2(t *testing.T) {

	var val Decimal2

	err := val.Scan([]byte("0.00"))

	fmt.Println(err, val)

}
