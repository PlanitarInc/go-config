package config

import "github.com/PlanitarInc/config/decoders"

type Flow struct {
	decoders []decoders.Decoder
	Config   interface{}
}

func NewFlow(defaults interface{}, ds ...decoders.Decoder) *Flow {
	return &Flow{
		decoders: ds,
		Config:   defaults,
	}
}

func (f Flow) Load() error {
	for _, d := range f.decoders {
		if err := d.Decode(f.Config); err != nil {
			return err
		}
	}
	return nil
}
