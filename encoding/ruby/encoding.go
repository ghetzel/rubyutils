package ruby

import (
	"bytes"
	"fmt"
	"github.com/fatih/structs"
	"github.com/ghetzel/go-stockutil/stringutil"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type encoderFunc func(e *encodeState, v reflect.Value) error // {}

type encodeState struct {
	bytes.Buffer
	indentEnabled bool
	indentLevel   int
	indent        []byte
	indentPrefix  []byte
}

func (self *encodeState) marshal(v interface{}) error {
	return self.reflectValue(reflect.ValueOf(v))
}

func (self *encodeState) reflectValue(v reflect.Value) error {
	return valueEncoder(v)(self, v)
}

func (self *encodeState) writeBytes(values ...[]byte) {
	if self.indentEnabled {
		if len(self.indentPrefix) > 0 {
			self.Write(self.indentPrefix)
		}

		for i := 0; i < self.indentLevel; i++ {
			self.Write(self.indent)
		}
	}

	for _, value := range values {
		self.Write(value)
	}
}

func (self *encodeState) writeStrings(values ...string) {
	byteset := make([][]byte, len(values))

	for i, value := range values {
		byteset[i] = []byte(value[:])
	}

	self.writeBytes(byteset...)
}

// retrieves a function capable of encoding the specific type of the input reflect.Value
func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}

	switch v.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32, reflect.Float64:
		return floatEncoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return structEncoder
	case reflect.Map:
		return mapEncoder
	case reflect.Slice:
		return sliceEncoder
	case reflect.Array:
		return arrayEncoder
	case reflect.Ptr:
		return ptrEncoder
	default:
		return unsupportedTypeEncoder
	}
}

// encode invalid values
func invalidValueEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(`nil`)
	return nil
}

func unsupportedTypeEncoder(e *encodeState, v reflect.Value) error {
	return fmt.Errorf("Unsupported type '%T', cannot encode", v.Interface())
}

// encode boolean values
func boolEncoder(e *encodeState, v reflect.Value) error {
	if v.Bool() {
		e.writeStrings(`true`)
	} else {
		e.writeStrings(`false`)
	}

	return nil
}

// encode integers
func intEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(strconv.FormatInt(v.Int(), 10))
	return nil
}

// encode unsigned integers
func uintEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(strconv.FormatUint(v.Uint(), 10))
	return nil
}

// encode floats
func floatEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(strconv.FormatFloat(v.Float(), 'f', -1, v.Type().Bits()))
	return nil
}

// encode strings (non-interpolated)
func stringEncoder(e *encodeState, v reflect.Value) error {
	str := v.String()

	// escape single-quotes
	str = strings.Replace(str, `'`, `\'`, -1)

	e.writeStrings(`'`, str, `'`)

	return nil
}

// encode a named key-value pair (in an output hash)
func keyValueEncoder(e *encodeState, key reflect.Value, value reflect.Value) error {
	keyEnc := &encodeState{}
	valEnc := &encodeState{}

	if err := valueEncoder(key)(keyEnc, key); err != nil {
		return err
	}

	if err := valueEncoder(value)(valEnc, value); err != nil {
		return err
	}

	e.writeBytes(keyEnc.Bytes(), []byte{' ', '=', '>', ' '}, valEnc.Bytes())
	return nil
}

// encode generic interfaces
func interfaceEncoder(e *encodeState, v reflect.Value) error {
	if v.IsNil() {
		e.writeStrings(`nil`)
		return nil
	}

	return e.reflectValue(v.Elem())
}

// encode structs
func structEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(`{`)
	e.indentLevel += 1

	defer func() {
		e.indentLevel -= 1
		e.writeStrings(`}`)
	}()

	structValue := structs.New(v.Interface())
	structFields := structValue.Fields()

	for i := 0; i < len(structFields); i++ {
		structField := structFields[i]

		var fieldName string
		var skip bool

		if structField.IsExported() {
			tagParts := strings.Split(structField.Tag(`ruby`), `,`)

			if len(tagParts) > 0 {
				fieldName = tagParts[0]
				tagParts = tagParts[1:]
			}

			for _, opt := range tagParts {
				switch opt {
				case `omitempty`:
					if structField.IsZero() {
						skip = true
					}
				}
			}
		} else {
			skip = true
		}

		// default struct field name to the field's name
		if fieldName == `` {
			fieldName = structField.Name()
		}

		// if we're not skipping this field, write it to the buffer
		if !skip {
			keyValueEncoder(e, reflect.ValueOf(structField.Name()), reflect.ValueOf(structField.Value()))

			// for all but the last element
			if i < (len(structFields) - 1) {
				e.writeStrings(`, `)

				if e.indentEnabled {
					e.writeStrings("\n")
				}
			}
		}
	}

	return nil
}

// encode maps with a best-attempt at deterministic ordering
func mapEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(`{`)
	e.indentLevel += 1

	defer func() {
		e.indentLevel -= 1
		e.writeStrings(`}`)
	}()

	keys := v.MapKeys()

	// this is a trick to provide ordered maps for keys types we can sort by
	strKeys := make([]string, 0)

	for _, key := range keys {
		if str, err := stringutil.ToString(key.Interface()); err == nil {
			strKeys = append(strKeys, str)
		}
	}

	if len(strKeys) == len(keys) {
		// only sort if we were able to stringify all the keys
		sort.Strings(strKeys)
	} else {
		// otherwise, make sure strKeys is the same size as keys because of what we do with it below
		strKeys = make([]string, len(keys))
	}

	// iterare over strKeys to get the correct index, but pull the actual key value from keys
	for i, _ := range strKeys {
		key := keys[i]
		value := v.MapIndex(key)

		if err := keyValueEncoder(e, key, value); err != nil {
			return err
		}

		// for all but the last element
		if i < (len(keys) - 1) {
			e.writeStrings(`, `)

			if e.indentEnabled {
				e.writeStrings("\n")
			}
		}
	}

	return nil
}

// encode slices
func sliceEncoder(e *encodeState, v reflect.Value) error {
	return nil
}

// encode arrays
func arrayEncoder(e *encodeState, v reflect.Value) error {
	return nil
}

// encode pointers to things
func ptrEncoder(e *encodeState, v reflect.Value) error {
	return nil
}
