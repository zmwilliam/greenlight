package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	value := fmt.Sprintf("%d mins", r)
	valueQuoted := strconv.Quote(value)
	return []byte(valueQuoted), nil
}
