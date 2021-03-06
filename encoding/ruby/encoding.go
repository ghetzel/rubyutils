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

type encodeStructField struct {
	Name  reflect.Value
	Value reflect.Value
}

func (self *encodeState) marshal(v interface{}) error {
	return self.reflectValue(reflect.ValueOf(v))
}

func (self *encodeState) reflectValue(v reflect.Value) error {
	return valueEncoder(v)(self, v)
}

func (self *encodeState) getIndentBytes() []byte {
	indentation := make([]byte, 0)

	if self.indentEnabled {
		if len(self.indentPrefix) > 0 {
			indentation = append(indentation, self.indentPrefix...)
		}

		for i := 0; i < self.indentLevel; i++ {
			indentation = append(indentation, self.indent...)
		}
	}

	return indentation
}

func (self *encodeState) writeBytesUnindented(values ...[]byte) {
	for _, value := range values {
		self.Write(value)
	}
}

func (self *encodeState) writeBytes(values ...[]byte) {
	self.Write(self.getIndentBytes())
	self.writeBytesUnindented(values...)
}

func (self *encodeState) writeStringsUnindented(values ...string) {
	byteset := make([][]byte, len(values))

	for i, value := range values {
		byteset[i] = []byte(value[:])
	}

	self.writeBytesUnindented(byteset...)
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
	case reflect.Interface, reflect.Ptr:
		return elementEncoder
	case reflect.Struct:
		return structEncoder
	case reflect.Map:
		return mapEncoder
	case reflect.Slice, reflect.Array:
		return arrayEncoder
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
	keyEnc := &encodeState{
		indentEnabled: e.indentEnabled,
		indentLevel:   (e.indentLevel - 1),
		indent:        e.indent,
		indentPrefix:  e.indentPrefix,
	}

	valEnc := &encodeState{
		indentEnabled: e.indentEnabled,
		indentLevel:   e.indentLevel,
		indent:        e.indent,
		indentPrefix:  e.indentPrefix,
	}

	if err := valueEncoder(key)(keyEnc, key); err != nil {
		return err
	}

	if err := valueEncoder(value)(valEnc, value); err != nil {
		return err
	}

	keyBytes := bytes.TrimPrefix(keyEnc.Bytes(), keyEnc.getIndentBytes())
	valBytes := bytes.TrimPrefix(valEnc.Bytes(), valEnc.getIndentBytes())

	var hashRocket []byte

	if e.indentEnabled {
		hashRocket = []byte{' ', '=', '>', ' '}
	} else {
		hashRocket = []byte{'=', '>'}
	}

	e.writeBytes(keyBytes, hashRocket, valBytes)
	return nil
}

// encode generic interfaces and pointers
func elementEncoder(e *encodeState, v reflect.Value) error {
	if v.IsNil() {
		e.writeStrings(`nil`)
		return nil
	}

	return e.reflectValue(v.Elem())
}

// encode structs
func structEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(`{`)

	structValue := structs.New(v.Interface())
	structFields := structValue.Fields()
	fieldsToWrite := make([]*encodeStructField, 0)

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

		// specifying a field name of "-" skips that field
		if fieldName == `-` {
			skip = true
		}

		// if we're not skipping this field, write it to the buffer
		if !skip {
			fieldsToWrite = append(fieldsToWrite, &encodeStructField{
				Name:  reflect.ValueOf(fieldName),
				Value: reflect.ValueOf(structField.Value()),
			})
		}
	}

	// only break and increase indentation if we have anything to write
	if len(fieldsToWrite) > 0 {
		if e.indentEnabled {
			e.writeStringsUnindented("\n")
		}

		e.indentLevel += 1
	}

	for i, field := range fieldsToWrite {
		keyValueEncoder(e, field.Name, field.Value)

		// for all but the last element
		if i < (len(fieldsToWrite) - 1) {
			e.writeStringsUnindented(`,`)

			// trail all but the last field with a space after the comma,
			// but only if we're not indenting (otherwise this space should be a line break)
			if !e.indentEnabled {
				e.writeStringsUnindented(` `)
			}
		}

		// add line break if we're indenting
		if e.indentEnabled {
			e.writeStringsUnindented("\n")
		}
	}

	if len(fieldsToWrite) > 0 {
		e.indentLevel -= 1
		e.writeStrings(`}`)
	} else {
		// if we didn't write anything, don't indent this (it's just an empty "{}")
		e.writeStringsUnindented(`}`)
	}

	return nil
}

// encode maps with a best-attempt at deterministic ordering by stringifying
// key names and outputing them in lexical order
func mapEncoder(e *encodeState, v reflect.Value) error {
	e.writeStrings(`{`)

	keys := v.MapKeys()

	// this is a trick to provide ordered maps for keys types we can sort by
	strKeys := make([]string, 0)
	sortKeyValues := make(map[string]reflect.Value)

	// attempt to convert the key value to a string and map that value to the key itself
	for _, key := range keys {
		if str, err := stringutil.ToString(key.Interface()); err == nil {
			strKeys = append(strKeys, str)
			sortKeyValues[str] = key
		}
	}

	if len(strKeys) == len(keys) {
		// only sort if we were able to stringify all the keys
		sort.Strings(strKeys)
	} else {
		// otherwise, make sure strKeys is the same size as keys because of what we do with it below
		strKeys = make([]string, len(keys))
	}

	// only break and increase indentation if we have anything to write
	if len(strKeys) > 0 {
		if e.indentEnabled {
			e.writeStringsUnindented("\n")
		}

		e.indentLevel += 1
	}

	// iterare over strKeys to get the correct index, but pull the actual key value from keys
	for i, sortKey := range strKeys {
		var key reflect.Value

		// attempt to retrieve the sorted key value, otherwise just get the next
		// one in sequence
		if k, ok := sortKeyValues[sortKey]; ok {
			key = k
		} else {
			key = keys[i]
		}

		// get the value at key
		value := v.MapIndex(key)

		// encode it
		if err := keyValueEncoder(e, key, value); err != nil {
			return err
		}

		// for all but the last element, add comma and (optionally) linebreak
		if i < (len(keys) - 1) {
			e.writeStringsUnindented(`,`)

			// trail all but the last field with a space after the comma,
			// but only if we're not indenting (otherwise this space should be a line break)
			if !e.indentEnabled {
				e.writeStringsUnindented(` `)
			}
		}

		// add line break if we're indenting
		if e.indentEnabled {
			e.writeStringsUnindented("\n")
		}
	}

	// only lower indentation and close if we had anything to write
	if len(strKeys) > 0 {
		e.indentLevel -= 1
		e.writeStrings(`}`)
	} else {
		// if we didn't write anything, don't indent this (it's just an empty "{}")
		e.writeStringsUnindented(`}`)
	}

	return nil
}

// encode arrays and slices
func arrayEncoder(e *encodeState, v reflect.Value) error {
	e.writeStringsUnindented(`[`)

	if e.indentEnabled {
		e.writeStringsUnindented("\n")
	}

	e.indentLevel += 1

	defer func() {
		e.indentLevel -= 1
		e.writeStringsUnindented(`]`)
	}()

	// iterare over strKeys to get the correct index, but pull the actual key value from keys
	for i := 0; i < v.Len(); i++ {
		value := v.Index(i)

		if err := valueEncoder(value)(e, value); err != nil {
			return err
		}

		// for all but the last element
		if i < (v.Len() - 1) {
			e.writeStringsUnindented(`,`)

			// trail all but the last field with a space after the comma,
			// but only if we're not indenting (otherwise this space should be a line break)
			if !e.indentEnabled {
				e.writeStringsUnindented(` `)
			}
		}

		// add line break if we're indenting
		if e.indentEnabled {
			e.writeStringsUnindented("\n")
		}
	}

	return nil
}
