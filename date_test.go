package validation_test

import (
	"testing"
	"time"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	tests := []struct {
		tag    string
		layout string
		value  interface{}
		err    string
	}{
		{"t1", time.ANSIC, "", ""},
		{"t2", time.ANSIC, "Wed Feb  4 21:00:57 2009", ""},
		{"t3", time.ANSIC, "Wed Feb  29 21:00:57 2009", "must be a valid date"},
		{"t4", "2006-01-02", "2009-11-12", ""},
		{"t5", "2006-01-02", "2009-11-12 21:00:57", "must be a valid date"},
		{"t6", "2006-01-02", "2009-1-12", "must be a valid date"},
		{"t7", "2006-01-02", "2009-01-12", ""},
		{"t8", "2006-01-02", "2009-01-32", "must be a valid date"},
		{"t9", "2006-01-02", 1, "must be either a string or byte slice"},
	}

	for _, test := range tests {
		r := validation.Date(test.layout)
		err := r.Validate(test.value)
		assertError(t, test.err, err, test.tag)
	}
}

func TestDateRule_Error(t *testing.T) {
	r := validation.Date(time.RFC3339)
	assert.Equal(t, "must be a valid date", r.Validate("0001-01-02T15:04:05Z07:00").Error())
	r2 := r.Min(time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC))
	assert.Equal(t, "the date is out of range", r2.Validate("1999-01-02T15:04:05Z").Error())
	r = r.Error("123")
	r = r.RangeError("456")
	assert.Equal(t, "123", r.Err.Message())
	assert.Equal(t, "456", r.RangeErr.Message())
}

func TestDateRule_ErrorObject(t *testing.T) {
	r := validation.Date(time.RFC3339)
	assert.Equal(t, "must be a valid date", r.Validate("0001-01-02T15:04:05Z07:00").Error())

	r = r.ErrorObject(validation.NewError("code", "abc"))

	assert.Equal(t, "code", r.Err.Code())
	assert.Equal(t, "abc", r.Validate("0001-01-02T15:04:05Z07:00").Error())

	r2 := r.Min(time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC))
	assert.Equal(t, "the date is out of range", r2.Validate("1999-01-02T15:04:05Z").Error())

	r = r.ErrorObject(validation.NewError("C", "def"))
	r = r.RangeErrorObject(validation.NewError("D", "123"))

	assert.Equal(t, "C", r.Err.Code())
	assert.Equal(t, "def", r.Err.Message())
	assert.Equal(t, "D", r.RangeErr.Code())
	assert.Equal(t, "123", r.RangeErr.Message())
}

func TestDateRule_MinMax(t *testing.T) {
	r := validation.Date(time.ANSIC)
	assert.True(t, r.Minimum.IsZero())
	assert.True(t, r.Maximum.IsZero())
	r = r.Min(time.Now())
	assert.False(t, r.Minimum.IsZero())
	assert.True(t, r.Maximum.IsZero())
	r = r.Max(time.Now())
	assert.False(t, r.Maximum.IsZero())

	r2 := validation.Date("2006-01-02").Min(time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC)).Max(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC))
	assert.Nil(t, r2.Validate("2010-01-02"))
	err := r2.Validate("1999-01-02")
	if assert.NotNil(t, err) {
		assert.Equal(t, "the date is out of range", err.Error())
	}
	err = r2.Validate("2021-01-02")
	if assert.NotNil(t, err) {
		assert.Equal(t, "the date is out of range", err.Error())
	}
}
