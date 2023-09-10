package validation

var (
	// ErrNil is the error that returns when a value is not nil.
	ErrNil = NewError("validation_nil", "must be blank")
	// ErrEmpty is the error that returns when a not nil value is not empty.
	ErrEmpty = NewError("validation_empty", "must be blank")
)

// Nil is a validation rule that checks if a value is nil.
// It is the opposite of NotNil rule
var Nil = absentRule{Condition: true, SkipNil: false}

// Empty checks if a not nil value is empty.
var Empty = absentRule{Condition: true, SkipNil: true}

type absentRule struct {
	Condition bool
	Err       Error
	SkipNil   bool
}

// Validate checks if the given value is valid or not.
func (r absentRule) Validate(value interface{}) error {
	if r.Condition {
		value, isNil := Indirect(value)
		if !r.SkipNil && !isNil || r.SkipNil && !isNil && !IsEmpty(value) {
			if r.Err != nil {
				return r.Err
			}
			if r.SkipNil {
				return ErrEmpty
			}
			return ErrNil
		}
	}
	return nil
}

// When sets the condition that determines if the validation should be performed.
func (r absentRule) When(condition bool) absentRule {
	r.Condition = condition
	return r
}

// Error sets the error message for the rule.
func (r absentRule) Error(message string) absentRule {
	if r.Err == nil {
		if r.SkipNil {
			r.Err = ErrEmpty
		} else {
			r.Err = ErrNil
		}
	}
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r absentRule) ErrorObject(err Error) absentRule {
	r.Err = err
	return r
}
