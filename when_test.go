package validation_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/prodadidb/go-validation"
)

func abcValidation(val string) bool {
	return val == "abc"
}

func TestWhen(t *testing.T) {
	abcRule := validation.NewStringRule(abcValidation, "wrong_abc")
	validateMeRule := validation.NewStringRule(validateMe, "wrong_me")

	tests := []struct {
		tag       string
		condition bool
		value     interface{}
		rules     []validation.Rule
		elseRules []validation.Rule
		err       string
	}{
		// True condition
		{"t1.1", true, nil, []validation.Rule{}, []validation.Rule{}, ""},
		{"t1.2", true, "", []validation.Rule{}, []validation.Rule{}, ""},
		{"t1.3", true, "", []validation.Rule{abcRule}, []validation.Rule{}, ""},
		{"t1.4", true, 12, []validation.Rule{validation.Required}, []validation.Rule{}, ""},
		{"t1.5", true, nil, []validation.Rule{validation.Required}, []validation.Rule{}, "cannot be blank"},
		{"t1.6", true, "123", []validation.Rule{abcRule}, []validation.Rule{}, "wrong_abc"},
		{"t1.7", true, "abc", []validation.Rule{abcRule}, []validation.Rule{}, ""},
		{"t1.8", true, "abc", []validation.Rule{abcRule, abcRule}, []validation.Rule{}, ""},
		{"t1.9", true, "abc", []validation.Rule{abcRule, validateMeRule}, []validation.Rule{}, "wrong_me"},
		{"t1.10", true, "me", []validation.Rule{abcRule, validateMeRule}, []validation.Rule{}, "wrong_abc"},
		{"t1.11", true, "me", []validation.Rule{}, []validation.Rule{abcRule}, ""},

		// False condition
		{"t2.1", false, "", []validation.Rule{}, []validation.Rule{}, ""},
		{"t2.2", false, "", []validation.Rule{abcRule}, []validation.Rule{}, ""},
		{"t2.3", false, "abc", []validation.Rule{abcRule}, []validation.Rule{}, ""},
		{"t2.4", false, "abc", []validation.Rule{abcRule, abcRule}, []validation.Rule{}, ""},
		{"t2.5", false, "abc", []validation.Rule{abcRule, validateMeRule}, []validation.Rule{}, ""},
		{"t2.6", false, "me", []validation.Rule{abcRule, validateMeRule}, []validation.Rule{}, ""},
		{"t2.7", false, "", []validation.Rule{abcRule, validateMeRule}, []validation.Rule{}, ""},
		{"t2.8", false, "me", []validation.Rule{}, []validation.Rule{abcRule, validateMeRule}, "wrong_abc"},
	}

	for _, test := range tests {
		err := validation.Validate(test.value, validation.When(test.condition, test.rules...).Else(test.elseRules...))
		assertError(t, test.err, err, test.tag)
	}
}

func TestWhenWithContext(t *testing.T) {
	type ctxKey int
	const (
		containsKey ctxKey = iota
	)
	rule := validation.WithContext(func(ctx context.Context, value interface{}) error {
		if !strings.Contains(value.(string), ctx.Value(containsKey).(string)) {
			return errors.New("unexpected value")
		}
		return nil
	})
	ctx1 := context.WithValue(context.Background(), containsKey, "abc")
	ctx2 := context.WithValue(context.Background(), containsKey, "xyz")

	tests := []struct {
		tag       string
		condition bool
		value     interface{}
		ctx       context.Context
		err       string
	}{
		// True condition
		{"t1.1", true, "abc", ctx1, ""},
		{"t1.2", true, "abc", ctx2, "unexpected value"},
		{"t1.3", true, "xyz", ctx1, "unexpected value"},
		{"t1.4", true, "xyz", ctx2, ""},

		// False condition
		{"t2.1", false, "abc", ctx1, ""},
		{"t2.2", false, "abc", ctx2, "unexpected value"},
		{"t2.3", false, "xyz", ctx1, "unexpected value"},
		{"t2.4", false, "xyz", ctx2, ""},
	}

	for _, test := range tests {
		err := validation.ValidateWithContext(test.ctx, test.value, validation.When(test.condition, rule).Else(rule))
		assertError(t, test.err, err, test.tag)
	}
}
