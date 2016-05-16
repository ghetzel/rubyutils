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
