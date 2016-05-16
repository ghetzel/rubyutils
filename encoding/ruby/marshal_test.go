package ruby

import (
	"testing"
)

func TestMarshalStruct(t *testing.T) {
	in := TestStruct{
		Name:        `test`,
		Count:       1,
		SkipIfZero:  2,
		SkipAlways:  true,
		notExported: 5,
	}

	shouldBe := `{'Name'=>'test', 'count'=>1, 'SkipIfZero'=>2}`

	if data, err := Marshal(in); err != nil {
		t.Fatal(err)
	} else {
		str := string(data[:])

		if str != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, str)
		} else {
			t.Log(str)
		}
	}
}

func TestMarshalIndentStruct(t *testing.T) {
	in := map[string]int{
		`second`: 2,
		`first`:  1,
		`third`:  3,
	}

	shouldBe := "{\n  'first' => 1,\n  'second' => 2,\n  'third' => 3\n}"

	if data, err := MarshalIndent(in, ``, `  `); err != nil {
		t.Fatal(err)
	} else {
		str := string(data[:])

		if str != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, str)
		} else {
			t.Log(str)
		}
	}
}
