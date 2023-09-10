package validation_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

func validateMe(s string) bool {
	return s == "me"
}

func TestNewStringRule(t *testing.T) {
	v := validation.NewStringRule(validateMe, "abc")

	assert.NotNil(t, v.StringValidate)
	assert.Equal(t, "", v.Err.Code())
	assert.Equal(t, "abc", v.Err.Message())
}

func TestNewStringRuleWithError(t *testing.T) {
	err := validation.NewError("C", "abc")
	v := validation.NewStringRuleWithError(validateMe, err)

	assert.NotNil(t, v.StringValidate)
	assert.Equal(t, err, v.Err)
	assert.Equal(t, "C", v.Err.Code())
	assert.Equal(t, "abc", v.Err.Message())
}

func TestStringRule_Error(t *testing.T) {
	err := validation.NewError("code", "abc")
	v := validation.NewStringRuleWithError(validateMe, err).Error("abc")
	assert.Equal(t, "code", v.Err.Code())
	assert.Equal(t, "abc", v.Err.Message())

	v2 := v.Error("correct")
	assert.Equal(t, "code", v.Err.Code())
	assert.Equal(t, "correct", v2.Err.Message())
	assert.Equal(t, "abc", v.Err.Message())
}

func TestStringValidator_Validate(t *testing.T) {
	v := validation.NewStringRule(validateMe, "wrong_rule").Error("wrong")

	value := "me"

	err := v.Validate(value)
	assert.Nil(t, err)

	err = v.Validate(&value)
	assert.Nil(t, err)

	value = ""

	err = v.Validate(value)
	assert.Nil(t, err)

	err = v.Validate(&value)
	assert.Nil(t, err)

	nullValue := sql.NullString{String: "me", Valid: true}
	err = v.Validate(nullValue)
	assert.Nil(t, err)

	nullValue = sql.NullString{String: "", Valid: true}
	err = v.Validate(nullValue)
	assert.Nil(t, err)

	var s *string
	err = v.Validate(s)
	assert.Nil(t, err)

	err = v.Validate("not me")
	if assert.NotNil(t, err) {
		assert.Equal(t, "wrong", err.Error())
	}

	err = v.Validate(100)
	if assert.NotNil(t, err) {
		assert.NotEqual(t, "wrong", err.Error())
	}

	v2 := v.Error("Wrong!")
	err = v2.Validate("not me")
	if assert.NotNil(t, err) {
		assert.Equal(t, "Wrong!", err.Error())
	}
}

func TestGetErrorFieldName(t *testing.T) {
	type A struct {
		T0 string
		T1 string `json:"t1"`
		T2 string `json:"t2,omitempty"`
		T3 string `json:",omitempty"`
		T4 string `json:"t4,x1,omitempty"` // nolint:staticcheck
	}
	tests := []struct {
		tag   string
		field string
		name  string
	}{
		{"t1", "T0", "T0"},
		{"t2", "T1", "t1"},
		{"t3", "T2", "t2"},
		{"t4", "T3", "T3"},
		{"t5", "T4", "t4"},
	}
	a := reflect.TypeOf(A{})
	for _, test := range tests {
		field, _ := a.FieldByName(test.field)
		assert.Equal(t, test.name, validation.GetErrorFieldName(&field), test.tag)
	}
}

func TestStringRule_ErrorObject(t *testing.T) {
	r := validation.NewStringRule(validateMe, "wrong_rule")

	err := validation.NewError("code", "abc")
	r = r.ErrorObject(err)

	assert.Equal(t, err, r.Err)
	assert.Equal(t, "code", r.Err.Code())
	assert.Equal(t, "abc", r.Err.Message())
}
