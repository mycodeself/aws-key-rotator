package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSliceExists(t *testing.T) {
	str := "my-str"
	s := []string{"my-str-1", str, "my-str-2"}

	res := StringInSlice(str, s)

	assert.True(t, res)
}

func TestStringInSliceNotExists(t *testing.T) {
	str := "str-not-exists"
	s := []string{"my-str-1", "my-str-2"}

	res := StringInSlice(str, s)

	assert.False(t, res)
}

func TestStringInSliceEmptyReturnFalse(t *testing.T) {
	str := "str-not-exists"
	s := []string{}

	res := StringInSlice(str, s)

	assert.False(t, res)
}