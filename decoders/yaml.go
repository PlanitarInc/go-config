package decoders

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlFileDecoder struct {
	filename string
}

func NewYamlFileDecoder(filename string) Decoder {
	return &yamlFileDecoder{
		filename: filename,
	}
}

func (d yamlFileDecoder) Decode(src interface{}) error {
	bs, err := ioutil.ReadFile(d.filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, src)
}
