package ruby

// import (
// 	"bytes"
// 	"reflect"
// 	"sync"
// )

func Marshal(v interface{}) ([]byte, error) {
	e := &encodeState{}
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

func MarshalIndent(v interface{}, prefix string, indent string) ([]byte, error) {
	e := &encodeState{
		indentEnabled: true,
		indentPrefix:  []byte(prefix[:]),
		indent:        []byte(indent[:]),
	}

	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}
