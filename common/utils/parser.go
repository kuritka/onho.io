// Package utils implements helpers.
package utils

import (
	"fmt"
	"strconv"
)

// ToFloat64 cast to float64.
func ToFloat64(i interface{}) (b float64) {
	switch v := i.(type) {
	case int:
		b = float64(v)
	case int8:
		b = float64(v)
	case int16:
		b = float64(v)
	case int32:
		b = float64(v)
	case int64:
		b = float64(v)
	case float32:
		b = float64(v)
	case float64:
		b = v
	case string:
		b, _ = strconv.ParseFloat(v, 64)
	default:
		b = 0
	}
	return
}

// ToFloat32 cast to float32.
func ToFloat32(i interface{}) float32 {
	return float32(ToFloat64(i))
}

// ToInt64 cast to int64.
func ToInt64(i interface{}) int64 {
	return int64(ToFloat64(i))
}

// ToInt32 cast to int32.
func ToInt32(i interface{}) int32 {
	return int32(ToFloat64(i))
}

// ToInt cast to int.
func ToInt(i interface{}) int {
	return int(ToFloat64(i))
}

// ToString cast to string.
func ToString(i interface{}) (s string) {
	if i == nil {
		return
	}

	switch i := i.(type) {
	case string:
		s = i
	default:
		s = fmt.Sprintf("%v", i)
	}

	return
}
