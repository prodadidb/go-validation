package validation

// ErrNotInInvalid is the error that returns when a value is in a list.
var ErrNotInInvalid = NewError("validation_not_in_invalid", "must not be in list")

// NotIn returns a validation rule that checks if a value is absent from the given list of values.
// Note that the value being checked and the possible range of values must be of the same type.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func NotIn(values ...interface{}) NotInRule {
	return NotInRule{
		Elements: values,
		Err:      ErrNotInInvalid,
	}
}

// NotInRule is a validation rule that checks if a value is absent from the given list of values.
type NotInRule struct {
	Elements []interface{}
	Err      Error
}

// Validate checks if the given value is valid or not.
func (r NotInRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	for _, e := range r.Elements {
		if e == value {
			return r.Err
		}
	}
	return nil
}

// Error sets the error message for the rule.
func (r NotInRule) Error(message string) NotInRule {
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r NotInRule) ErrorObject(err Error) NotInRule {
	r.Err = err
	return r
}
