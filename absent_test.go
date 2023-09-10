package validation_test

import (
	"testing"
	"time"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	s1 := "123"
	s2 := ""
	var time1 time.Time
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", 123, "must be blank"},
		{"t2", "", "must be blank"},
		{"t3", &s1, "must be blank"},
		{"t4", &s2, "must be blank"},
		{"t5", nil, ""},
		{"t6", time1, "must be blank"},
	}

	for _, test := range tests {
		r := validation.Nil
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestEmpty(t *testing.T) {
	s1 := "123"
	s2 := ""
	time1 := time.Now()
	var time2 time.Time
	tests := []struct {
		tag   string
		value interface{}
		err   string
	}{
		{"t1", 123, "must be blank"},
		{"t2", "", ""},
		{"t3", &s1, "must be blank"},
		{"t4", &s2, ""},
		{"t5", nil, ""},
		{"t6", time1, "must be blank"},
		{"t7", time2, ""},
	}

	for _, test := range tests {
		r := validation.Empty
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestAbsentRule_When(t *testing.T) {
	r := validation.Nil.When(false)
	err := validation.Validate(42, r)
	assert.Nil(t, err)

	r = validation.Nil.When(true)
	err = validation.Validate(42, r)
	assert.Equal(t, validation.ErrNil, err)
}

func Test_absentRule_Error(t *testing.T) {
	r := validation.Nil
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.False(t, r.SkipNil)
	r2 := r.Error("123")
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.False(t, r.SkipNil)
	assert.Equal(t, "123", r2.Err.Message())
	assert.False(t, r2.SkipNil)

	r = validation.Empty
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.True(t, r.SkipNil)
	r2 = r.Error("123")
	assert.Equal(t, "must be blank", r.Validate("42").Error())
	assert.True(t, r.SkipNil)
	assert.Equal(t, "123", r2.Err.Message())
	assert.True(t, r2.SkipNil)
}

func TestAbsentRule_Error(t *testing.T) {
	r := validation.Nil

	err := validation.NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.Err)
	assert.Equal(t, err.Code(), r.Err.Code())
	assert.Equal(t, err.Message(), r.Err.Message())
	assert.NotEqual(t, err, validation.Nil.Err)
}
