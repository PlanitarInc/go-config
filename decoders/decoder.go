package decoders

import (
	"io/ioutil"
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

type Unmarshaller interface {
	Unmarshall([]byte, interface{}) error
}

type fileunmarshaller struct {
	filename string
	u        Unmarshaller
}

func (f fileunmarshaller) Decode(dst interface{}) error {
	bs, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return err
	}
	return f.u.Unmarshall(bs, dst)
}

func NewFileUnmarshaller(filename string, u Unmarshaller) Decoder {
	return &fileunmarshaller{
		filename: filename,
		u:        u,
	}
}
