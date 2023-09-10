package validation_test

import (
	"testing"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestNotIn(t *testing.T) {
	v := 1
	var v2 *int
	var tests = []struct {
		tag    string
		values []interface{}
		value  interface{}
		err    string
	}{
		{"t0", []interface{}{1, 2}, 0, ""},
		{"t1", []interface{}{1, 2}, 1, "must not be in list"},
		{"t2", []interface{}{1, 2}, 2, "must not be in list"},
		{"t3", []interface{}{1, 2}, 3, ""},
		{"t4", []interface{}{}, 3, ""},
		{"t5", []interface{}{1, 2}, "1", ""},
		{"t6", []interface{}{1, 2}, &v, "must not be in list"},
		{"t7", []interface{}{1, 2}, v2, ""},
	}

	for _, test := range tests {
		r := validation.NotIn(test.values...)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func Test_NotInRule_Error(t *testing.T) {
	r := validation.NotIn(1, 2, 3)
	assert.Equal(t, "must not be in list", r.Validate(1).Error())
	r = r.Error("123")
	assert.Equal(t, "123", r.Err.Message())
}

func TestNotInRule_ErrorObject(t *testing.T) {
	r := validation.NotIn(1, 2, 3)

	err := validation.NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.Err)
	assert.Equal(t, err.Code(), r.Err.Code())
	assert.Equal(t, err.Message(), r.Err.Message())
}
