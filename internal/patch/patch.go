package patch

import (
	"bytes"
	"encoding/json"
)

type Field[T any] struct {
	Set bool
	Val *T
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	f.Set = true

	if bytes.Equal(data, []byte("null")) {
		f.Val = nil
		return nil
	}

	var v T

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	f.Val = &v
	return nil
}

func (f *Field[T]) IsSet() bool {
	return f.Set
}

func (f Field[T]) HasValue() bool { //----------------Проверка на наличие значение в переменной Val и Set==true
	return f.Set && f.Val != nil
}
