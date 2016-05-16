package ruby

import (
	"strings"
	"testing"
)

func TestEncodeIndentedSimpleMapStrInt(t *testing.T) {
	e := &encodeState{
		indentEnabled: true,
		indent:        []byte{' ', ' '},
	}

	in := map[string]int{
		`second`: 2,
		`first`:  1,
		`third`:  3,
	}

	shouldBe := "{\n  'first' => 1,\n  'second' => 2,\n  'third' => 3\n}"

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected:\n\"%s\",\ngot: \"%s\"", shouldBe, s)
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeIndentedSimpleMapStrMixed(t *testing.T) {
	e := &encodeState{
		indentEnabled: true,
		indent:        []byte{' ', ' '},
	}

	in := map[string]interface{}{
		`third`:  9.6,
		`second`: 4,
		`first`:  true,
	}

	shouldBe := "{\n  'first' => true,\n  'second' => 4,\n  'third' => 9.6\n}"

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected:\n\"%s\",\ngot: \"%s\"", shouldBe, s)
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeIndentedSliceInt(t *testing.T) {
	e := &encodeState{
		indentEnabled: true,
		indent:        []byte{' ', ' '},
	}

	in := []int{1, 2, 3}
	shouldBe := "[\n  1,\n  2,\n  3\n]"

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected:\n\"%s\",\ngot: \"%s\"", shouldBe, strings.Replace(s, " ", ".", -1))
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeIndentedSliceString(t *testing.T) {
	e := &encodeState{
		indentEnabled: true,
		indent:        []byte{' ', ' '},
	}

	in := []string{`one`, `two`, `three`}
	shouldBe := "[\n  'one',\n  'two',\n  'three'\n]"

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected:\n\"%s\",\ngot: \"%s\"", shouldBe, strings.Replace(s, " ", ".", -1))
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeIndentedStruct(t *testing.T) {
	e := &encodeState{
		indentEnabled: true,
		indent:        []byte{' ', ' '},
	}

	in := TestStruct{
		Name:        `test`,
		Count:       1,
		SkipIfZero:  2,
		SkipAlways:  true,
		notExported: 5,
	}

	shouldBe := "{\n  'Name' => 'test',\n  'count' => 1,\n  'SkipIfZero' => 2\n}"

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, strings.Replace(s, " ", ".", -1))
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeIndentedStructComplex(t *testing.T) {
	e := &encodeState{
		indentEnabled: true,
		indent:        []byte{' ', ' '},
	}

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
			`prop-4`: map[string]interface{}{},
		},
	}

	shouldBe := "{\n  'Name' => 'test',\n  'Data' => {\n    'Key' => 'first-level',\n    'Value' => {\n      'Key' => 'second-level'\n    }\n  },\n  'Properties' => {\n    'prop-1' => true,\n    'prop-2' => 4,\n    'prop-3' => {\n      'subprop-1' => {\n        'Key' => 'subnest-1',\n        'Value' => {}\n      }\n    },\n    'prop-4' => {}\n  }\n}"

	if err := e.marshal(in); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, strings.Replace(s, " ", ".", -1))
		} else {
			t.Log(s)
		}
	}
}
