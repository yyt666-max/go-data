package encoding

import (
	"encoding"
	"encoding/json"
)

type Marshaller interface {
	json.Marshaler
	encoding.BinaryMarshaler
}
type Unmarshaler interface {
	encoding.BinaryUnmarshaler
	json.Unmarshaler
}

type MarshalUnmarshaler interface {
	Marshaller
	Unmarshaler
}

var (
	_ Unmarshaler = (*jsonUnmarshaler[any])(nil)
	_ Marshaller  = (*jsonUnmarshaler[any])(nil)
)

type jsonUnmarshaler[T any] struct {
	obj *T
}

func Json[T any](t *T) MarshalUnmarshaler {
	return &jsonUnmarshaler[T]{
		obj: t,
	}
}
func JsonUnmarshaler[T any](t *T) Unmarshaler {
	return &jsonUnmarshaler[T]{
		obj: t,
	}
}
func JsonMarshaller[T any](t *T) Unmarshaler {
	return &jsonUnmarshaler[T]{
		obj: t,
	}
}

func (e *jsonUnmarshaler[T]) UnmarshalBinary(data []byte) error {
	return e.UnmarshalJSON(data)
}

func (e *jsonUnmarshaler[T]) MarshalBinary() (data []byte, err error) {
	return e.MarshalJSON()
}

func (e *jsonUnmarshaler[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.obj)
}

func (e *jsonUnmarshaler[T]) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &e.obj)
}
