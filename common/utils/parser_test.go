package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testsNumbers = []struct {
	n        interface{}
	expected float64
}{
	{20, 20},
	{"", 0},
	{0, 0},
	{int(-20), -20},
	{int8(-20), -20},
	{int16(-20), -20},
	{int32(-20), -20},
	{int64(-20), -20},
	{float32(-20), -20},
	{float64(-20), -20},
	{true, 0},
	{false, 0},
	{"-20", -20},
	{"0", 0},
	{"12 12", 0},
	{"    12", 0},
	{"12    ", 0},
	{"    12    ", 0},
	{nil, 0},
}

func TestToFloat64(t *testing.T) {
	for i := range testsNumbers {
		actual := ToFloat64(testsNumbers[i].n)
		assert.Equal(t, testsNumbers[i].expected, actual)
	}
}

func TestToFloat32(t *testing.T) {
	for i := range testsNumbers {
		actual := ToFloat32(testsNumbers[i].n)
		assert.Equal(t, float32(testsNumbers[i].expected), actual)
	}
}

func TestToInt64(t *testing.T) {
	for i := range testsNumbers {
		actual := ToInt64(testsNumbers[i].n)
		assert.Equal(t, int64(testsNumbers[i].expected), actual)

	}
}

func TestToInt32(t *testing.T) {
	for i := range testsNumbers {
		actual := ToInt32(testsNumbers[i].n)
		assert.Equal(t, int32(testsNumbers[i].expected), actual)
	}
}

func TestToInt(t *testing.T) {
	for i := range testsNumbers {
		actual := ToInt(testsNumbers[i].n)
		assert.Equal(t, int(testsNumbers[i].expected), actual)
	}
}

var testStrings = []struct {
	n        interface{}
	expected string
}{
	{0, "0"},
	{1, "1"},
	{nil, ""},
	{"1", "1"},
	{"1", "1"},
}

func TestToString(t *testing.T) {
	for i := range testStrings {
		actual := ToString(testStrings[i].n)
		assert.Equal(t, testStrings[i].expected, actual)
	}
}
