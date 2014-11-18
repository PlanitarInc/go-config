package decoders

import "encoding/json"

type jsonUnmarshaller struct{}

func (u jsonUnmarshaller) Unmarshall(bs []byte, dst interface{}) error {
	return json.Unmarshal(bs, dst)
}

func NewJsonFileDecoder(filename string) Decoder {
	return NewFileUnmarshaller(filename, &jsonUnmarshaller{})
}
