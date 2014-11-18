package config

import (
	"errors"
	"os"
	"testing"

	"github.com/PlanitarInc/config/decoders"
	. "github.com/onsi/gomega"
)

func TestEmptyFlow(t *testing.T) {
	RegisterTestingT(t)

	type Cfg struct {
		Number int
		Flag   bool
		Str1   string
		Str2   string
	}

	def := Cfg{Number: -123, Flag: true, Str1: "qwe", Str2: "asd"}

	actual := def
	f := NewFlow(&actual)
	Ω(f.LoadFailIfError()).ShouldNot(HaveOccurred())
	Ω(f.Config).Should(Equal(&def))
}

func TestEnvFlow(t *testing.T) {
	RegisterTestingT(t)

	type Cfg struct {
		Number int
		Flag   bool
		Str1   string
		Str2   string
	}

	def := Cfg{Number: -123, Flag: true, Str1: "qwe", Str2: "asd"}

	actual := def
	f := NewFlow(&actual, decoders.NewEnvDecoder(""))
	os.Setenv("NUMBER", "-987")
	os.Setenv("STR2", "asd-987")
	exp := def
	exp.Number = -987
	exp.Str2 = "asd-987"
	Ω(f.LoadFailIfError()).ShouldNot(HaveOccurred())
	Ω(f.Config).Should(Equal(&exp))
}

func TestYamlFileFlow(t *testing.T) {
	RegisterTestingT(t)

	type Cfg struct {
		Number int
		Flag   bool
		Str1   string
		Str2   string
	}

	def := Cfg{Number: -123, Flag: true, Str1: "qwe", Str2: "asd"}

	actual := def
	f := NewFlow(&actual, decoders.NewYamlFileDecoder("test.yml"))
	exp := def
	exp.Number = 987654321
	exp.Str1 = "aasd"
	Ω(f.LoadFailIfError()).ShouldNot(HaveOccurred())
	Ω(f.Config).Should(Equal(&exp))
}

func TestJsonFileFlow(t *testing.T) {
	RegisterTestingT(t)

	type Cfg struct {
		Number int
		Flag   bool
		Str1   string
		Str2   string
	}

	def := Cfg{Number: -123, Flag: true, Str1: "qwe", Str2: "asd"}

	actual := def
	f := NewFlow(&actual, decoders.NewJsonFileDecoder("test.json"))
	exp := def
	exp.Number = -17
	exp.Flag = false
	Ω(f.LoadFailIfError()).ShouldNot(HaveOccurred())
	Ω(f.Config).Should(Equal(&exp))
}

func TestFlowOrder(t *testing.T) {
	RegisterTestingT(t)

	type Cfg struct {
		Number int
		Flag   bool
		Str1   string
		Str2   string
	}

	def := Cfg{Number: -123, Flag: true, Str1: "qwe", Str2: "asd"}

	actual := def
	ds := []decoders.Decoder{
		// YAML config file override the default settings
		decoders.NewYamlFileDecoder("test.yml"),
		// JSON config file override the YAML settings
		decoders.NewJsonFileDecoder("test.json"),
		// Env vars override the JSON settings
		decoders.NewEnvDecoder(""),
		// Command line arguments override the Env var settings
		// XXX
	}
	f := NewFlow(&actual, ds...)

	os.Setenv("NUMBER", "-987")
	os.Setenv("STR2", "asd-987")

	exp := def
	exp.Number = -987
	exp.Flag = false
	exp.Str1 = "aasd"
	exp.Str2 = "asd-987"
	Ω(f.LoadFailIfError()).ShouldNot(HaveOccurred())
	Ω(f.Config).Should(Equal(&exp))
}

type FailDecoder struct {
	err error
}

func (d FailDecoder) Decode(dst interface{}) error {
	return d.err
}

type CountDecoder struct {
	cnt int
}

func (d *CountDecoder) Decode(dst interface{}) error {
	d.cnt++
	return nil
}

func TestErrorCollection(t *testing.T) {
	RegisterTestingT(t)

	e1 := errors.New("1")
	e2 := errors.New("2")
	e3 := errors.New("3")
	act := struct{}{}
	var f *Flow
	var cd *CountDecoder

	cd = &CountDecoder{}
	f = NewFlow(&act, cd)
	Ω(f.Load()).Should(Equal([]error{}))
	Ω(cd.cnt).Should(Equal(1))
	Ω(f.LoadFailIfError()).ShouldNot(HaveOccurred())
	Ω(cd.cnt).Should(Equal(2))

	cd = &CountDecoder{}
	f = NewFlow(&act, FailDecoder{e1}, cd)
	Ω(f.Load()).Should(Equal([]error{e1}))
	Ω(cd.cnt).Should(Equal(1))
	Ω(f.LoadFailIfError()).Should(Equal(e1))
	Ω(cd.cnt).Should(Equal(1))

	cd = &CountDecoder{}
	f = NewFlow(&act, cd, &FailDecoder{e1}, cd, &FailDecoder{e2}, cd)
	Ω(f.Load()).Should(Equal([]error{e1, e2}))
	Ω(cd.cnt).Should(Equal(3))
	Ω(f.LoadFailIfError()).Should(Equal(e1))
	Ω(cd.cnt).Should(Equal(4))

	cd = &CountDecoder{}
	f = NewFlow(&act, cd, FailDecoder{e1}, FailDecoder{e3}, FailDecoder{e2},
		cd)
	Ω(f.Load()).Should(Equal([]error{e1, e3, e2}))
	Ω(cd.cnt).Should(Equal(2))
	Ω(f.LoadFailIfError()).Should(Equal(e1))
	Ω(cd.cnt).Should(Equal(3))
}
