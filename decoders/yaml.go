package decoders

import "gopkg.in/yaml.v2"

type yamlUnmarshaller struct{}

func (u yamlUnmarshaller) Unmarshall(bs []byte, dst interface{}) error {
	return yaml.Unmarshal(bs, dst)
}

func NewYamlFileDecoder(filename string) Decoder {
	return NewFileUnmarshaller(filename, &yamlUnmarshaller{})
}
