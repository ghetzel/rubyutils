package ruby

import (
	"testing"
)

func TestEncodeBool(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(true); err != nil {
		t.Fatal(err)
	}

	if s := e.String(); s != `true` {
		t.Fatalf("Expected 'true', got '%s'", s)
	} else {
		t.Log(s)
	}
}

func TestEncodeInt(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(5); err != nil {
		t.Fatal(err)
	}

	if s := e.String(); s != `5` {
		t.Fatalf("Expected '5', got '%s'", s)
	} else {
		t.Log(s)
	}
}

func TestEncodeIntNegative(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(int64(-9223372036854775808)); err != nil {
		t.Fatal(err)
	}

	if s := e.String(); s != `-9223372036854775808` {
		t.Fatalf("Expected '-9223372036854775808', got '%s'", s)
	} else {
		t.Log(s)
	}
}

func TestEncodeIntBig(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(9223372036854775807); err != nil {
		t.Fatal(err)
	}

	if s := e.String(); s != `9223372036854775807` {
		t.Fatalf("Expected '9223372036854775807', got '%s'", s)
	} else {
		t.Log(s)
	}
}

func TestEncodeUint(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(uint(8)); err != nil {
		t.Fatal(err)
	}

	if s := e.String(); s != `8` {
		t.Fatalf("Expected '8', got '%s'", s)
	} else {
		t.Log(s)
	}
}

func TestEncodeUintBig(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(uint64(18446744073709551615)); err != nil {
		t.Fatal(err)
	}

	if s := e.String(); s != `18446744073709551615` {
		t.Fatalf("Expected '18446744073709551615', got '%s'", s)
	} else {
		t.Log(s)
	}
}

func TestEncodeFloat32(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(float32(0)); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != `0` {
			t.Fatalf("Expected '0', got '%s'", s)
		} else {
			t.Log(s)
		}
	}

	e.Reset()

	if err := e.marshal(float32(0.1)); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != `0.1` {
			t.Fatalf("Expected '0.1', got '%s'", s)
		} else {
			t.Log(s)
		}
	}

	e.Reset()

	// last digit should truncate, round the 6 up to 7
	if err := e.marshal(float32(3.14159265)); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != `3.1415927` {
			t.Fatalf("Expected '3.1415927', got '%s'", s)
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeFloat64(t *testing.T) {
	e := &encodeState{}

	if err := e.marshal(float64(0)); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != `0` {
			t.Fatalf("Expected '0', got '%s'", s)
		} else {
			t.Log(s)
		}
	}

	e.Reset()

	if err := e.marshal(float64(0.1)); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != `0.1` {
			t.Fatalf("Expected '0.1', got '%s'", s)
		} else {
			t.Log(s)
		}
	}

	e.Reset()

	// last digit should truncate
	if err := e.marshal(float64(3.1415926535897932)); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != `3.141592653589793` {
			t.Fatalf("Expected '3.141592653589793', got '%s'", s)
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeString(t *testing.T) {
	e := &encodeState{}

	shouldBe := `'test'`

	if err := e.marshal(`test`); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, s)
		} else {
			t.Log(s)
		}
	}

	e.Reset()

	shouldBe := `'test\'s test'`

	if err := e.marshal(`test's test`); err != nil {
		t.Fatal(err)
	} else {
		if s := e.String(); s != shouldBe {
			t.Fatalf("Expected \"%s\", got \"%s\"", shouldBe, s)
		} else {
			t.Log(s)
		}
	}
}

func TestEncodeSimpleMapStrInt(t *testing.T) {
	e := &encodeState{}

	in := map[string]int{
		`second`: 2,
		`first`:  1,
		`third`:  3,
	}

	shouldBe := `{'first' => 1, 'second' => 2, 'third' => 3}`

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

func TestEncodeSimpleMapStrMixed(t *testing.T) {
	e := &encodeState{}

	in := map[string]interface{}{
		`third`:  9.6,
		`second`: 4,
		`first`:  true,
	}

	shouldBe := `{'first' => true, 'second' => 4, 'third' => 9.6}`

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

func TestEncodeSliceInt(t *testing.T) {
	e := &encodeState{}

	in := []int{1, 2, 3}
	shouldBe := `[1, 2, 3]`

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

func TestEncodeSliceString(t *testing.T) {
	e := &encodeState{}

	in := []string{`one`, `two`, `three`}
	shouldBe := `['one', 'two', 'three']`

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
