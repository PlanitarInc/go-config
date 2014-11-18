package decoders

import (
	"reflect"

	"github.com/PlanitarInc/config/reflectx"
)

type Decoder interface {
	Decode(dst interface{}) error
}

type KVStore interface {
	DecodeKey(key string, dst interface{}) error
	Tagname() string
	MapFunc() func(string) string
	ReduceFunc() func(string, string) string
}

type kvwrapper struct {
	mapper *reflectx.Mapper
	store  KVStore
}

func (d kvwrapper) Decode(dst interface{}) error {
	m := d.mapper.FieldMap(reflect.ValueOf(dst))
	for key, field := range m {
		pf := field.Addr().Interface()
		if err := d.store.DecodeKey(key, pf); err != nil {
			return err
		}
	}
	return nil
}

func KVWrapper(s KVStore) Decoder {
	m := reflectx.NewMapperFunc(s.Tagname(), s.MapFunc())
	m.SetReduceFunc(s.ReduceFunc())
	return &kvwrapper{mapper: m, store: s}
}
