package decoders

import (
	"encoding/json"
	"io/ioutil"
)

type jsonFileDecoder struct {
	filename string
}

func NewJsonFileDecoder(filename string) Decoder {
	return &jsonFileDecoder{
		filename: filename,
	}
}

func (d jsonFileDecoder) Decode(src interface{}) error {
	bs, err := ioutil.ReadFile(d.filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, src)
}
