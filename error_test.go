package validation_test

import (
	"errors"
	"testing"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

func TestNewInternalError(t *testing.T) {
	err := validation.NewInternalError(errors.New("abc"))
	if assert.NotNil(t, err.InternalError()) {
		assert.Equal(t, "abc", err.InternalError().Error())
	}
}

func TestErrors_Error(t *testing.T) {
	errs := validation.Errors{
		"B": errors.New("B1"),
		"C": errors.New("C1"),
		"A": errors.New("A1"),
	}
	assert.Equal(t, "A: A1; B: B1; C: C1.", errs.Error())

	errs = validation.Errors{
		"B": errors.New("B1"),
	}
	assert.Equal(t, "B: B1.", errs.Error())

	errs = validation.Errors{}
	assert.Equal(t, "", errs.Error())
}

func TestErrors_MarshalMessage(t *testing.T) {
	errs := validation.Errors{
		"A": errors.New("A1"),
		"B": validation.Errors{
			"2": errors.New("B1"),
		},
	}
	errsJSON, err := errs.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, "{\"A\":\"A1\",\"B\":{\"2\":\"B1\"}}", string(errsJSON))
}

func TestErrors_Filter(t *testing.T) {
	errs := validation.Errors{
		"B": errors.New("B1"),
		"C": nil,
		"A": errors.New("A1"),
	}
	err := errs.Filter()
	assert.Equal(t, 2, len(errs))
	if assert.NotNil(t, err) {
		assert.Equal(t, "A: A1; B: B1.", err.Error())
	}

	errs = validation.Errors{}
	assert.Nil(t, errs.Filter())

	errs = validation.Errors{
		"B": nil,
		"C": nil,
	}

	assert.Nil(t, errs.Filter())
}

func TestErrorObject_SetCode(t *testing.T) {
	err := validation.NewError("A", "msg").(validation.ErrorObject)

	assert.Equal(t, err.ErrCode, "A")
	assert.Equal(t, err.Code(), "A")

	err = err.SetCode("B").(validation.ErrorObject)
	assert.Equal(t, "B", err.ErrCode)
}

func TestErrorObject_Code(t *testing.T) {
	err := validation.NewError("A", "msg").(validation.ErrorObject)

	assert.Equal(t, err.Code(), "A")
}

func TestErrorObject_SetMessage(t *testing.T) {
	err := validation.NewError("code", "A").(validation.ErrorObject)

	assert.Equal(t, err.ErrMessage, "A")
	assert.Equal(t, err.Message(), "A")

	err = err.SetMessage("abc").(validation.ErrorObject)
	assert.Equal(t, err.ErrMessage, "abc")
	assert.Equal(t, err.Message(), "abc")
}

func TestErrorObject_Message(t *testing.T) {
	err := validation.NewError("code", "A").(validation.ErrorObject)

	assert.Equal(t, err.ErrMessage, "A")
	assert.Equal(t, err.Message(), "A")
}

func TestErrorObject_Params(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "AA": "val2"}

	err := validation.NewError("code", "A").(validation.ErrorObject)
	err = err.SetParams(p).(validation.ErrorObject)
	err = err.SetMessage("B").(validation.ErrorObject)

	assert.Equal(t, err.ErrParams, p)
	assert.Equal(t, err.Params(), p)
}

func TestErrorObject_AddParam2(t *testing.T) {
	p := map[string]interface{}{"key": "val"}
	err := validation.NewError("code", "A").(validation.ErrorObject)
	err = err.AddParam("key", "val").(validation.ErrorObject)

	assert.Equal(t, err.ErrParams, p)
	assert.Equal(t, err.Params(), p)
}

func TestErrorObject_AddParam(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "B": "val2"}

	err := validation.NewError("code", "A").(validation.ErrorObject)
	err = err.SetParams(p).(validation.ErrorObject)
	err = err.AddParam("C", "val3").(validation.ErrorObject)

	p["C"] = "val3"

	assert.Equal(t, err.ErrParams, p)
	assert.Equal(t, err.Params(), p)
}

func TestError_Code(t *testing.T) {
	err := validation.NewError("A", "msg")

	assert.Equal(t, err.Code(), "A")
}

func TestError_SetMessage(t *testing.T) {
	err := validation.NewError("code", "A")

	assert.Equal(t, err.Message(), "A")

	err = err.SetMessage("abc")
	assert.Equal(t, err.Message(), "abc")
}

func TestError_Message(t *testing.T) {
	err := validation.NewError("code", "A")

	assert.Equal(t, err.Message(), "A")
}

func TestError_Params(t *testing.T) {
	p := map[string]interface{}{"A": "val1", "AA": "val2"}

	err := validation.NewError("code", "A")
	err = err.SetParams(p)
	err = err.SetMessage("B")

	assert.Equal(t, err.Params(), p)
}

func TestValidationError(t *testing.T) {
	params := map[string]interface{}{
		"A": "B",
	}

	err := validation.NewError("code", "msg")
	err = err.SetParams(params)

	assert.Equal(t, err.Code(), "code")
	assert.Equal(t, err.Message(), "msg")
	assert.Equal(t, err.Params(), params)

	params = map[string]interface{}{"min": 1}
	err = err.SetParams(params)

	assert.Equal(t, err.Params(), params)
}
