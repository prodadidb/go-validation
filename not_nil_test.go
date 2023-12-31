package validation_test

import (
	"testing"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

type MyInterface interface {
	Hello()
}

func TestNotNil(t *testing.T) {
	var v1 []int
	var v2 map[string]int
	var v3 *int
	var v4 interface{}
	var v5 MyInterface
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", v1, "is required"},
		{"t2", v2, "is required"},
		{"t3", v3, "is required"},
		{"t4", v4, "is required"},
		{"t5", v5, "is required"},
		{"t6", "", ""},
		{"t7", 0, ""},
	}

	for _, test := range tests {
		r := validation.NotNil
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func Test_notNilRule_Error(t *testing.T) {
	r := validation.NotNil
	assert.Equal(t, "is required", r.Validate(nil).Error())
	r2 := r.Error("123")
	assert.Equal(t, "is required", r.Validate(nil).Error())
	assert.Equal(t, "123", r2.Err.Message())
}

func TestNotNilRule_ErrorObject(t *testing.T) {
	r := validation.NotNil

	err := validation.NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.Err)
	assert.Equal(t, err.Code(), r.Err.Code())
	assert.Equal(t, err.Message(), r.Err.Message())
	assert.NotEqual(t, err, validation.NotNil.Err)
}
