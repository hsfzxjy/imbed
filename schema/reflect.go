package schema

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sync"

	"github.com/hsfzxjy/imbed/util"
)

var ratType = reflect.TypeOf((*big.Rat)(nil))

type Store struct {
	mu sync.RWMutex
	m  map[reflect.Type]GenericSchema
}

func NewStore() *Store {
	return &Store{
		m: map[reflect.Type]GenericSchema{},
	}
}

func Register[T any](store *Store) (GenericSchema, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	var dummy T
	typ := reflect.TypeOf(dummy)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expect a struct, got %s", typ.String())
	} else if s, ok := store.m[typ]; ok {
		return s, nil
	}
	nFields := typ.NumField()
	fields := make([]*_StructField, 0, nFields)
	for i := 0; i < nFields; i++ {
		f := typ.Field(i)
		include := false
		if f.Tag == "imbed" {
			include = true
		} else if _, ok := f.Tag.Lookup("imbed"); ok {
			include = true
		}
		if !include {
			continue
		}
		elemSchema, err := store.schemaFor(f.Type)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", f.Name, err)
		}
		field := &_StructField{
			name:       f.Name,
			offset:     f.Offset,
			elemSchema: elemSchema,
		}
		fields = append(fields, field)
	}
	schema := new_Toplevel(new_Struct[T](typ.Name(), fields))
	store.m[typ] = schema
	return schema, nil
}

func RegisterMust[T any](store *Store) GenericSchema {
	return util.Unwrap(Register[T](store))
}

var ErrNotRegistered = errors.New("no schema registered")

func (s *Store) Get(v any) (GenericSchema, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Store.Get() expect *T, got %T", v)
	}
	return s.get(typ.Elem())
}

func (s *Store) get(typ reflect.Type) (GenericSchema, error) {
	if sch, ok := s.m[typ]; ok {
		return sch, nil
	} else {
		return nil, fmt.Errorf("%w for %s", ErrNotRegistered, typ.String())
	}
}

func (s *Store) Lookup(v any) (GenericSchema, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Ptr {
		return nil, false
	}
	sch, ok := s.m[typ.Elem()]
	return sch, ok
}

func (s *Store) schemaFor(typ reflect.Type) (genericSchema, error) {
	switch typ.Kind() {
	case reflect.Int64:
		return new_Int(optional[int64]{}), nil
	case reflect.Ptr:
		if typ == ratType {
			return new_Rat(optional[*big.Rat]{}), nil
		}
	case reflect.Bool:
		return new_Bool(optional[bool]{}), nil
	case reflect.String:
		return new_String(optional[string]{}), nil
	case reflect.Struct:
		return s.get(typ)
	}
	return nil, fmt.Errorf("cannot create schema for %s", typ.String())
}
