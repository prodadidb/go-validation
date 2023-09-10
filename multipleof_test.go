package validation_test

import (
	"testing"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestMultipleOf(t *testing.T) {
	r := validation.MultipleOf(10)
	assert.Equal(t, "must be multiple of 10", r.Validate(11).Error())
	assert.Equal(t, nil, r.Validate(20))
	assert.Equal(t, "cannot convert float32 to int64", r.Validate(float32(20)).Error())

	r2 := validation.MultipleOf("some string ....")
	assert.Equal(t, "type not supported: string", r2.Validate(10).Error())

	r3 := validation.MultipleOf(uint(10))
	assert.Equal(t, "must be multiple of 10", r3.Validate(uint(11)).Error())
	assert.Equal(t, nil, r3.Validate(uint(20)))
	assert.Equal(t, "cannot convert float32 to uint64", r3.Validate(float32(20)).Error())

}

func Test_MultipleOf_Error(t *testing.T) {
	r := validation.MultipleOf(10)
	assert.Equal(t, "must be multiple of 10", r.Validate(3).Error())

	r = r.Error("some error string ...")
	assert.Equal(t, "some error string ...", r.Err.Message())
}

func TestMultipleOfRule_ErrorObject(t *testing.T) {
	r := validation.MultipleOf(10)
	err := validation.NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.Err)
	assert.Equal(t, err.Code(), r.Err.Code())
	assert.Equal(t, err.Message(), r.Err.Message())
}
