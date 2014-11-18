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

func (f Flow) load(failOnError bool) []error {
	errs := []error{}
	for _, d := range f.decoders {
		if err := d.Decode(f.Config); err != nil {
			if failOnError {
				return []error{err}
			}
			errs = append(errs, err)
		}
	}
	return errs
}

func (f Flow) LoadFailIfError() error {
	errs := f.load(true)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (f Flow) Load() []error {
	return f.load(false)
}
