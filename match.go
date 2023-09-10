package validation

import (
	"regexp"
)

// ErrMatchInvalid is the error that returns in case of invalid format.
var ErrMatchInvalid = NewError("validation_match_invalid", "must be in a valid format")

// Match returns a validation rule that checks if a value matches the specified regular expression.
// This rule should only be used for validating strings and byte slices, or a validation error will be reported.
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func Match(re *regexp.Regexp) MatchRule {
	return MatchRule{
		Re:  re,
		Err: ErrMatchInvalid,
	}
}

// MatchRule is a validation rule that checks if a value matches the specified regular expression.
type MatchRule struct {
	Re  *regexp.Regexp
	Err Error
}

// Validate checks if the given value is valid or not.
func (r MatchRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil {
		return nil
	}

	isString, str, isBytes, bs := StringOrBytes(value)
	if isString && (str == "" || r.Re.MatchString(str)) {
		return nil
	} else if isBytes && (len(bs) == 0 || r.Re.Match(bs)) {
		return nil
	}
	return r.Err
}

// Error sets the error message for the rule.
func (r MatchRule) Error(message string) MatchRule {
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r MatchRule) ErrorObject(err Error) MatchRule {
	r.Err = err
	return r
}
