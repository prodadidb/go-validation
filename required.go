package validation

var (
	// ErrRequired is the error that returns when a value is required.
	ErrRequired = NewError("validation_required", "cannot be blank")
	// ErrNilOrNotEmpty is the error that returns when a value is not nil and is empty.
	ErrNilOrNotEmpty = NewError("validation_nil_or_not_empty_required", "cannot be blank")
)

// Required is a validation rule that checks if a value is not empty.
// A value is considered not empty if
// - integer, float: not zero
// - bool: true
// - string, array, slice, map: len() > 0
// - interface, pointer: not nil and the referenced value is not empty
// - any other types
var Required = RequiredRule{SkipNil: false, Condition: true}

// NilOrNotEmpty checks if a value is a nil pointer or a value that is not empty.
// NilOrNotEmpty differs from Required in that it treats a nil pointer as valid.
var NilOrNotEmpty = RequiredRule{SkipNil: true, Condition: true}

// RequiredRule is a rule that checks if a value is not empty.
type RequiredRule struct {
	Condition bool
	SkipNil   bool
	Err       Error
}

// Validate checks if the given value is valid or not.
func (r RequiredRule) Validate(value interface{}) error {
	if r.Condition {
		value, isNil := Indirect(value)
		if r.SkipNil && !isNil && IsEmpty(value) || !r.SkipNil && (isNil || IsEmpty(value)) {
			if r.Err != nil {
				return r.Err
			}
			if r.SkipNil {
				return ErrNilOrNotEmpty
			}
			return ErrRequired
		}
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r RequiredRule) When(condition bool) RequiredRule {
	r.Condition = condition
	return r
}

// Error sets the error message for the rule.
func (r RequiredRule) Error(message string) RequiredRule {
	if r.Err == nil {
		if r.SkipNil {
			r.Err = ErrNilOrNotEmpty
		} else {
			r.Err = ErrRequired
		}
	}
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r RequiredRule) ErrorObject(err Error) RequiredRule {
	r.Err = err
	return r
}
