package validation

import (
	"time"
)

var (
	// ErrDateInvalid is the error that returns in case of an invalid date.
	ErrDateInvalid = NewError("validation_date_invalid", "must be a valid date")
	// ErrDateOutOfRange is the error that returns in case of an invalid date.
	ErrDateOutOfRange = NewError("validation_date_out_of_range", "the date is out of range")
)

// DateRule is a validation rule that validates date/time string values.
type DateRule struct {
	Layout           string
	Minimum, Maximum time.Time
	Err, RangeErr    Error
}

// Date returns a validation rule that checks if a string value is in a format that can be parsed into a date.
// The format of the date should be specified as the Layout parameter which accepts the same value as that for time.Parse.
// For example,
//
//	validation.Date(time.ANSIC)
//	validation.Date("02 Jan 06 15:04 MST")
//	validation.Date("2006-01-02")
//
// By calling Min() and/or Max(), you can let the Date rule to check if a parsed date value is within
// the specified date range.
//
// An empty value is considered valid. Use the Required rule to make sure a value is not empty.
func Date(layout string) DateRule {
	return DateRule{
		Layout:   layout,
		Err:      ErrDateInvalid,
		RangeErr: ErrDateOutOfRange,
	}
}

// Error sets the error message that is used when the value being validated is not a valid date.
func (r DateRule) Error(message string) DateRule {
	r.Err = r.Err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct that is used when the value being validated is not a valid date..
func (r DateRule) ErrorObject(err Error) DateRule {
	r.Err = err
	return r
}

// RangeError sets the error message that is used when the value being validated is out of the specified Min/Max date range.
func (r DateRule) RangeError(message string) DateRule {
	r.RangeErr = r.RangeErr.SetMessage(message)
	return r
}

// RangeErrorObject sets the error struct that is used when the value being validated is out of the specified Min/Max date range.
func (r DateRule) RangeErrorObject(err Error) DateRule {
	r.RangeErr = err
	return r
}

// Min sets the minimum date range. A zero value means skipping the minimum range validation.
func (r DateRule) Min(min time.Time) DateRule {
	r.Minimum = min
	return r
}

// Max sets the maximum date range. A zero value means skipping the maximum range validation.
func (r DateRule) Max(max time.Time) DateRule {
	r.Maximum = max
	return r
}

// Validate checks if the given value is a valid date.
func (r DateRule) Validate(value interface{}) error {
	value, isNil := Indirect(value)
	if isNil || IsEmpty(value) {
		return nil
	}

	str, err := EnsureString(value)
	if err != nil {
		return err
	}

	date, err := time.Parse(r.Layout, str)
	if err != nil {
		return r.Err
	}

	if !r.Minimum.IsZero() && r.Minimum.After(date) || !r.Maximum.IsZero() && date.After(r.Maximum) {
		return r.RangeErr
	}

	return nil
}
