package validation

type stringValidator func(string) bool

// StringRule is a rule that checks a string variable using a specified stringValidator.
type StringRule struct {
	StringValidate stringValidator
	Err            Error
}

// NewStringRule creates a new validation rule using a function that takes a string value and returns a bool.
// The rule returned will use the function to check if a given string or byte slice is valid or not.
// An empty value is considered to be valid. Please use the Required rule to make sure a value is not empty.
func NewStringRule(validator stringValidator, message string) StringRule {
	return StringRule{
		StringValidate: validator,
		Err:            NewError("", message),
	}
}

// NewStringRuleWithError creates a new validation rule using a function that takes a string value and returns a bool.
// The rule returned will use the function to check if a given string or byte slice is valid or not.
// An empty value is considered to be valid. Please use the Required rule to make sure a value is not empty.
func NewStringRuleWithError(validator stringValidator, err Error) StringRule {
	return StringRule{
		StringValidate: validator,
		Err:            err,
	}
}

// Error sets the error message for the rule.
func (r StringRule) Error(message string) StringRule {
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r StringRule) ErrorObject(err Error) StringRule {
	r.Err = err
	return r
}

// Validate checks if the given value is valid or not.
func (r StringRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	str, err := EnsureString(value)
	if err != nil {
		return err
	}

	if r.StringValidate(str) {
		return nil
	}

	return r.Err
}
