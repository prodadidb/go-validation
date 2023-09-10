package validation_test

import (
	"testing"
	"time"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	s1 := "123"
	s2 := ""
	var time1 time.Time
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", 123, ""},
		{"t2", "", "cannot be blank"},
		{"t3", &s1, ""},
		{"t4", &s2, "cannot be blank"},
		{"t5", nil, "cannot be blank"},
		{"t6", time1, "cannot be blank"},
	}

	for _, test := range tests {
		r := validation.Required
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestRequiredRule_When(t *testing.T) {
	r := validation.Required.When(false)
	err := validation.Validate(nil, r)
	assert.Nil(t, err)

	r = validation.Required.When(true)
	err = validation.Validate(nil, r)
	assert.Equal(t, validation.ErrRequired, err)
}

func TestNilOrNotEmpty(t *testing.T) {
	s1 := "123"
	s2 := ""
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", 123, ""},
		{"t2", "", "cannot be blank"},
		{"t3", &s1, ""},
		{"t4", &s2, "cannot be blank"},
		{"t5", nil, ""},
	}

	for _, test := range tests {
		r := validation.NilOrNotEmpty
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func Test_requiredRule_Error(t *testing.T) {
	r := validation.Required
	assert.Equal(t, "cannot be blank", r.Validate(nil).Error())
	assert.False(t, r.SkipNil)
	r2 := r.Error("123")
	assert.Equal(t, "cannot be blank", r.Validate(nil).Error())
	assert.False(t, r.SkipNil)
	assert.Equal(t, "123", r2.Err.Message())
	assert.False(t, r2.SkipNil)

	r = validation.NilOrNotEmpty
	assert.Equal(t, "cannot be blank", r.Validate("").Error())
	assert.True(t, r.SkipNil)
	r2 = r.Error("123")
	assert.Equal(t, "cannot be blank", r.Validate("").Error())
	assert.True(t, r.SkipNil)
	assert.Equal(t, "123", r2.Err.Message())
	assert.True(t, r2.SkipNil)
}

func TestRequiredRule_Error(t *testing.T) {
	r := validation.Required

	err := validation.NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.Err)
	assert.Equal(t, err.Code(), r.Err.Code())
	assert.Equal(t, err.Message(), r.Err.Message())
	assert.NotEqual(t, err, validation.Required.Err)
}
