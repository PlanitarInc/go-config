package decoders

import (
	"reflect"

	"github.com/PlanitarInc/config/reflectx"
)

type structStore struct {
	fieldMap   map[string]reflect.Value
	dstTagname string
	reduceFunc func(string, string) string
}

func (s structStore) DecodeKey(key string, dst interface{}) error {
	if v, ok := s.fieldMap[key]; ok {
		reflect.Indirect(reflect.ValueOf(dst)).Set(v)
	}
	return nil
}

func (s structStore) Tagname() string {
	return s.dstTagname
}

func (s structStore) MapFunc() func(string) string {
	return nil
}

func (s structStore) ReduceFunc() func(string, string) string {
	return s.reduceFunc
}

func NewStructDecoder(src interface{}, srctag, dsttag string) Decoder {
	reduceFunc := reflectx.DelimiterKeyReducer(".")
	m := reflectx.NewMapper(srctag)
	m.SetReduceFunc(reduceFunc)
	return KVWrapper(&structStore{
		fieldMap:   m.FieldMap(reflect.ValueOf(src)),
		dstTagname: dsttag,
		reduceFunc: reduceFunc,
	})
}
