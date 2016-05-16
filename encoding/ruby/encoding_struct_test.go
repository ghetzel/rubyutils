package ruby

import (
	"testing"
)

type TestStruct struct {
	Name        string
	Count       int  `ruby:"count"`
	SkipIfZero  int  `ruby:",omitempty"`
	SkipAlways  bool `ruby:"-"`
	notExported int
}

type TestStructComplexSubNested struct {
	Key   string      `ruby:",omitempty"`
	Value interface{} `ruby:",omitempty"`
}

type TestStructComplexNested struct {
	Key   string
	Value TestStructComplexSubNested
}

type TestStructComplex struct {
	Name         string
	Data         TestStructComplexNested
	OptionalData *TestStructComplexNested `ruby:",omitempty"`
	Properties   map[string]interface{}
}

func TestEncodeStruct(t *testing.T) {
	e := &encodeState{}

	in := TestStruct{
		Name:        `test`,
		Count:       1,
		SkipIfZero:  2,
		SkipAlways:  true,
		notExported: 5,
	}

	shouldBe := `{'Name'=>'test', 'count'=>1, 'SkipIfZero'=>2}`

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, s)
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeStructComplex(t *testing.T) {
	e := &encodeState{}

	in := TestStructComplex{
		Name: `test`,
		Data: TestStructComplexNested{
			Key: `first-level`,
			Value: TestStructComplexSubNested{
				Key: `second-level`,
			},
		},
		Properties: map[string]interface{}{
			`prop-1`: true,
			`prop-2`: 4,
			`prop-3`: map[string]interface{}{
				`subprop-1`: &TestStructComplexNested{
					Key: `subnest-1`,
				},
			},
		},
	}

	shouldBe := `{'Name'=>'test', 'Data'=>{'Key'=>'first-level', 'Value'=>{'Key'=>'second-level'}}, 'Properties'=>{'prop-1'=>true, 'prop-2'=>4, 'prop-3'=>{'subprop-1'=>{'Key'=>'subnest-1', 'Value'=>{}}}}}`

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, s)
		} else {
			t.Log(s)
		}
	}
}
