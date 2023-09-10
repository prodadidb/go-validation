package validation

// ErrNotNilRequired is the error that returns when a value is Nil.
var ErrNotNilRequired = NewError("validation_not_nil_required", "is required")

// NotNil is a validation rule that checks if a value is not nil.
// NotNil only handles types including interface, pointer, slice, and map.
// All other types are considered valid.
var NotNil = notNilRule{}

type notNilRule struct {
	Err Error
}

// Validate checks if the given value is valid or not.
func (r notNilRule) Validate(value interface{}) error {
	_, isNil := Indirect(value)
	if isNil {
		if r.Err != nil {
			return r.Err
		}
		return ErrNotNilRequired
	}
	return nil
}

// Error sets the error message for the rule.
func (r notNilRule) Error(message string) notNilRule {
	if r.Err == nil {
		r.Err = ErrNotNilRequired
	}
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r notNilRule) ErrorObject(err Error) notNilRule {
	r.Err = err
	return r
}
