package dtoregistry

import (
	"errors"
	"fmt"
	"reflect"
)

type Registry struct {
	m map[reflect.Type]Converter
}

type Converter func(any) (any, error)

func NewRegistry() *Registry {
	return &Registry{
		m: make(map[reflect.Type]Converter),
	}
}

func (r *Registry) Register[I, O any](fn func(I) O) {
	r.m[reflect.TypeFor[I]()] = func(val any) (any, error) {
		in, ok := val.(I)
		if !ok {
			return nil, fmt.Errorf("expected %v, received %T", reflect.TypeFor[I](), val)
		}
		out := fn(in)

		return out, nil
	}
}

var ErrNoConverter = errors.New("converter not registered")

func (r *Registry) Match(v any) (any, error) {
	t := reflect.TypeOf(v)
	if f, ok := r.m[t]; ok {
		return f(v)
	}
	if t != nil && t.Kind() == reflect.Slice {
		converter, ok := r.m[t.Elem()]
		if !ok {
			return nil, fmt.Errorf("%w for %v", ErrNoConverter, t.Elem())
		}
		rv := reflect.ValueOf(v)
		out := make([]any, rv.Len())
		for i := range out {
			o, err := converter(rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			out[i] = o
		}
		return out, nil
	}
	return nil, fmt.Errorf("%w for %T", ErrNoConverter, v)
}
