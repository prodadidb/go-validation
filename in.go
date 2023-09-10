package validation

import (
	"reflect"
)

// ErrInInvalid is the error that returns in case of an invalid value for "in" rule.
var ErrInInvalid = NewError("validation_in_invalid", "must be a valid value")

// In returns a validation rule that checks if a value can be found in the given list of values.
// reflect.DeepEqual() will be used to determine if two values are equal.
// For more details please refer to https://golang.org/pkg/reflect/#DeepEqual
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func In(values ...interface{}) InRule {
	return InRule{
		Elements: values,
		Err:      ErrInInvalid,
	}
}

// InRule is a validation rule that validates if a value can be found in the given list of values.
type InRule struct {
	Elements []interface{}
	Err      Error
}

// Validate checks if the given value is valid or not.
func (r InRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	for _, e := range r.Elements {
		if reflect.DeepEqual(e, value) {
			return nil
		}
	}

	return r.Err
}

// Error sets the error message for the rule.
func (r InRule) Error(message string) InRule {
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r InRule) ErrorObject(err Error) InRule {
	r.Err = err
	return r
}
