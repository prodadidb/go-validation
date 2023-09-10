package validation_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/prodadidb/go-validation"
	"github.com/stretchr/testify/assert"
)

type Struct1 struct {
	Field1 int
	Field2 *int
	Field3 []int
	Field4 [4]int
	field5 int
	Struct2
	S1               *Struct2
	S2               Struct2
	JSONField        int `json:"some_json_field"`
	JSONIgnoredField int `json:"-"`
}

type Struct2 struct {
	Field21 string
	Field22 string
}

type Struct3 struct {
	*Struct2
	S1 string
}

func TestFindStructField(t *testing.T) {
	var s1 Struct1
	v1 := reflect.ValueOf(&s1).Elem()
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Field1)))
	assert.Nil(t, validation.FindStructField(v1, reflect.ValueOf(s1.Field2)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Field2)))
	assert.Nil(t, validation.FindStructField(v1, reflect.ValueOf(s1.Field3)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Field3)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Field4)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.field5)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Struct2)))
	assert.Nil(t, validation.FindStructField(v1, reflect.ValueOf(s1.S1)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.S1)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Field21)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Field22)))
	assert.NotNil(t, validation.FindStructField(v1, reflect.ValueOf(&s1.Struct2.Field22)))
	s2 := reflect.ValueOf(&s1.Struct2).Elem()
	assert.NotNil(t, validation.FindStructField(s2, reflect.ValueOf(&s1.Field21)))
	assert.NotNil(t, validation.FindStructField(s2, reflect.ValueOf(&s1.Struct2.Field21)))
	assert.NotNil(t, validation.FindStructField(s2, reflect.ValueOf(&s1.Struct2.Field22)))
	s3 := Struct3{
		Struct2: &Struct2{},
	}
	v3 := reflect.ValueOf(&s3).Elem()
	assert.NotNil(t, validation.FindStructField(v3, reflect.ValueOf(&s3.Struct2)))
	assert.NotNil(t, validation.FindStructField(v3, reflect.ValueOf(&s3.Field21)))
}

func TestValidateStruct(t *testing.T) {
	var m0 *Model1
	m1 := Model1{A: "abc", B: "xyz", c: "abc", G: "xyz", H: []string{"abc", "abc"}, I: map[string]string{"foo": "abc"}}
	m2 := Model1{E: String123("xyz")}
	m3 := Model2{}
	m4 := Model2{M3: Model3{A: "abc"}, Model3: Model3{A: "abc"}}
	m5 := Model2{Model3: Model3{A: "internal"}}
	tests := []struct {
		tag   string
		model interface{}
		rules []*validation.FieldRules
		err   string
	}{
		// empty rules
		{"t1.1", &m1, []*validation.FieldRules{}, ""},
		{"t1.2", &m1, []*validation.FieldRules{validation.Field(&m1.A), validation.Field(&m1.B)}, ""},
		// normal rules
		{"t2.1", &m1, []*validation.FieldRules{validation.Field(&m1.A, &validateAbc{}), validation.Field(&m1.B, &validateXyz{})}, ""},
		{"t2.2", &m1, []*validation.FieldRules{validation.Field(&m1.A, &validateXyz{}), validation.Field(&m1.B, &validateAbc{})}, "A: error xyz; B: error abc."},
		{"t2.3", &m1, []*validation.FieldRules{validation.Field(&m1.A, &validateXyz{}), validation.Field(&m1.c, &validateXyz{})}, "A: error xyz; c: error xyz."},
		{"t2.4", &m1, []*validation.FieldRules{validation.Field(&m1.D, validation.Length(0, 5))}, ""},
		{"t2.5", &m1, []*validation.FieldRules{validation.Field(&m1.F, validation.Length(0, 5))}, ""},
		{"t2.6", &m1, []*validation.FieldRules{validation.Field(&m1.H, validation.Each(&validateAbc{})), validation.Field(&m1.I, validation.Each(&validateAbc{}))}, ""},
		{"t2.7", &m1, []*validation.FieldRules{validation.Field(&m1.H, validation.Each(&validateXyz{})), validation.Field(&m1.I, validation.Each(&validateXyz{}))}, "H: (0: error xyz; 1: error xyz.); I: (foo: error xyz.)."},
		// non-struct pointer
		{"t3.1", m1, []*validation.FieldRules{}, validation.ErrStructPointer.Error()},
		{"t3.2", nil, []*validation.FieldRules{}, validation.ErrStructPointer.Error()},
		{"t3.3", m0, []*validation.FieldRules{}, ""},
		{"t3.4", &m0, []*validation.FieldRules{}, validation.ErrStructPointer.Error()},
		// invalid field spec
		{"t4.1", &m1, []*validation.FieldRules{validation.Field(m1)}, validation.ErrFieldPointer(0).Error()},
		{"t4.2", &m1, []*validation.FieldRules{validation.Field(&m1)}, validation.ErrFieldNotFound(0).Error()},
		// struct tag
		{"t5.1", &m1, []*validation.FieldRules{validation.Field(&m1.G, &validateAbc{})}, "g: error abc."},
		// validatable field
		{"t6.1", &m2, []*validation.FieldRules{validation.Field(&m2.E)}, "E: error 123."},
		{"t6.2", &m2, []*validation.FieldRules{validation.Field(&m2.E, validation.Skip)}, ""},
		{"t6.3", &m2, []*validation.FieldRules{validation.Field(&m2.E, validation.Skip.When(true))}, ""},
		{"t6.4", &m2, []*validation.FieldRules{validation.Field(&m2.E, validation.Skip.When(false))}, "E: error 123."},
		// Required, NotNil
		{"t7.1", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.Required)}, "F: cannot be blank."},
		{"t7.2", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.NotNil)}, "F: is required."},
		{"t7.3", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.Skip, validation.Required)}, ""},
		{"t7.4", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.Skip, validation.NotNil)}, ""},
		{"t7.5", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.Skip.When(true), validation.Required)}, ""},
		{"t7.6", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.Skip.When(true), validation.NotNil)}, ""},
		{"t7.7", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.Skip.When(false), validation.Required)}, "F: cannot be blank."},
		{"t7.8", &m2, []*validation.FieldRules{validation.Field(&m2.F, validation.Skip.When(false), validation.NotNil)}, "F: is required."},
		// embedded structs
		{"t8.1", &m3, []*validation.FieldRules{validation.Field(&m3.M3, validation.Skip)}, ""},
		{"t8.2", &m3, []*validation.FieldRules{validation.Field(&m3.M3)}, "M3: (A: error abc.)."},
		{"t8.3", &m3, []*validation.FieldRules{validation.Field(&m3.Model3, validation.Skip)}, ""},
		{"t8.4", &m3, []*validation.FieldRules{validation.Field(&m3.Model3)}, "A: error abc."},
		{"t8.5", &m4, []*validation.FieldRules{validation.Field(&m4.M3)}, ""},
		{"t8.6", &m4, []*validation.FieldRules{validation.Field(&m4.Model3)}, ""},
		{"t8.7", &m3, []*validation.FieldRules{validation.Field(&m3.A, validation.Required), validation.Field(&m3.B, validation.Required)}, "A: cannot be blank; B: cannot be blank."},
		{"t8.8", &m3, []*validation.FieldRules{validation.Field(&m4.A, validation.Required)}, "field #0 cannot be found in the struct"},
		// internal error
		{"t9.1", &m5, []*validation.FieldRules{validation.Field(&m5.A, &validateAbc{}), validation.Field(&m5.B, validation.Required), validation.Field(&m5.A, &validateInternalError{})}, "error internal"},
	}
	for _, test := range tests {
		err1 := validation.ValidateStruct(test.model, test.rules...)
		err2 := validation.ValidateStructWithContext(context.Background(), test.model, test.rules...)
		assertError(t, test.err, err1, test.tag)
		assertError(t, test.err, err2, test.tag)
	}

	// embedded struct
	err := validation.Validate(&m3)
	assert.EqualError(t, err, "A: error abc.")

	a := struct {
		Name  string
		Value string
	}{"name", "demo"}
	err = validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Value, validation.Required, validation.Length(5, 10)),
	)
	assert.EqualError(t, err, "Value: the length must be between 5 and 10.")
}

func TestValidateStructWithContext(t *testing.T) {
	m1 := Model1{A: "abc", B: "xyz", c: "abc", G: "xyz"}
	m2 := Model2{Model3: Model3{A: "internal"}}
	m3 := Model5{}
	tests := []struct {
		tag   string
		model interface{}
		rules []*validation.FieldRules
		err   string
	}{
		// normal rules
		{"t1.1", &m1, []*validation.FieldRules{validation.Field(&m1.A, &validateContextAbc{}), validation.Field(&m1.B, &validateContextXyz{})}, ""},
		{"t1.2", &m1, []*validation.FieldRules{validation.Field(&m1.A, &validateContextXyz{}), validation.Field(&m1.B, &validateContextAbc{})}, "A: error xyz; B: error abc."},
		{"t1.3", &m1, []*validation.FieldRules{validation.Field(&m1.A, &validateContextXyz{}), validation.Field(&m1.c, &validateContextXyz{})}, "A: error xyz; c: error xyz."},
		{"t1.4", &m1, []*validation.FieldRules{validation.Field(&m1.G, &validateContextAbc{})}, "g: error abc."},
		// skip rule
		{"t2.1", &m1, []*validation.FieldRules{validation.Field(&m1.G, validation.Skip, &validateContextAbc{})}, ""},
		{"t2.2", &m1, []*validation.FieldRules{validation.Field(&m1.G, &validateContextAbc{}, validation.Skip)}, "g: error abc."},
		// internal error
		{"t3.1", &m2, []*validation.FieldRules{validation.Field(&m2.A, &validateContextAbc{}), validation.Field(&m2.B, validation.Required), validation.Field(&m2.A, &validateInternalError{})}, "error internal"},
	}
	for _, test := range tests {
		err := validation.ValidateStructWithContext(context.Background(), test.model, test.rules...)
		assertError(t, test.err, err, test.tag)
	}

	//embedded struct
	err := validation.ValidateWithContext(context.Background(), &m3)
	if assert.NotNil(t, err) {
		assert.Equal(t, "A: error abc.", err.Error())
	}

	a := struct {
		Name  string
		Value string
	}{"name", "demo"}
	err = validation.ValidateStructWithContext(context.Background(), &a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Value, validation.Required, validation.Length(5, 10)),
	)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Value: the length must be between 5 and 10.", err.Error())
	}
}

func Test_GetErrorFieldName(t *testing.T) {
	var s1 Struct1
	v1 := reflect.ValueOf(&s1).Elem()

	sf1 := validation.FindStructField(v1, reflect.ValueOf(&s1.Field1))
	assert.NotNil(t, sf1)
	assert.Equal(t, "Field1", validation.GetErrorFieldName(sf1))

	jsonField := validation.FindStructField(v1, reflect.ValueOf(&s1.JSONField))
	assert.NotNil(t, jsonField)
	assert.Equal(t, "some_json_field", validation.GetErrorFieldName(jsonField))

	jsonIgnoredField := validation.FindStructField(v1, reflect.ValueOf(&s1.JSONIgnoredField))
	assert.NotNil(t, jsonIgnoredField)
	assert.Equal(t, "JSONIgnoredField", validation.GetErrorFieldName(jsonIgnoredField))
}
