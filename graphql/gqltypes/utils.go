package gqltypes

import "strconv"

func parseUint8FromStringOrPanic(s string) uint8 {
	val, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		panic(err)
	}
	return uint8(val)
}
