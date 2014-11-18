package decoders

import (
	"os"
	"strings"

	"github.com/PlanitarInc/go-config/reflectx"
	"gopkg.in/yaml.v2"
)

type envDecoder struct {
	tagname string
}

func (s envDecoder) DecodeKey(key string, dst interface{}) error {
	if val := os.Getenv(key); val != "" {
		// Let the yaml decoder do the hard work
		return yaml.Unmarshal([]byte(val), dst)
	}
	return nil
}

func (s envDecoder) Tagname() string {
	return s.tagname
}

func (s envDecoder) MapFunc() func(string) string {
	return strings.ToUpper
}

func (s envDecoder) ReduceFunc() func(string, string) string {
	return reflectx.DelimiterKeyReducer("_")
}

func NewEnvDecoder(tagname string) Decoder {
	return KVWrapper(&envDecoder{tagname})
}
