package config

import (
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
	Ω(f.Load()).ShouldNot(HaveOccurred())
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
	Ω(f.Load()).ShouldNot(HaveOccurred())
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
	Ω(f.Load()).ShouldNot(HaveOccurred())
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
	Ω(f.Load()).ShouldNot(HaveOccurred())
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
	Ω(f.Load()).ShouldNot(HaveOccurred())
	Ω(f.Config).Should(Equal(&exp))
}
