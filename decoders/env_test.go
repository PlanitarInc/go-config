package decoders

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func TestEnvDecodeKey(t *testing.T) {
	RegisterTestingT(t)

	e := envDecoder{}

	dStr := ""
	dInt := 0
	dFlt := 0.0

	os.Setenv("KEY", "asd")
	Ω(e.DecodeKey("KEY", &dStr)).ShouldNot(HaveOccurred())
	Ω(dStr).Should(Equal("asd"))

	os.Setenv("KEY", "123")
	Ω(e.DecodeKey("KEY", &dStr)).ShouldNot(HaveOccurred())
	Ω(dStr).Should(Equal("123"))

	Ω(e.DecodeKey("KEY", &dInt)).ShouldNot(HaveOccurred())
	Ω(dInt).Should(Equal(123))

	Ω(e.DecodeKey("KEY", &dFlt)).ShouldNot(HaveOccurred())
	Ω(dFlt).Should(Equal(123.0))

	os.Setenv("KEY", "")
}

func TestEnvDecoderBasic(t *testing.T) {
	RegisterTestingT(t)

	dst := struct {
		A int
		B string
		C float64
	}{}

	os.Setenv("A", "123")
	os.Setenv("B", "str")
	os.Setenv("C", "12.2")

	d := NewEnvDecoder("")
	d.Decode(&dst)

	Ω(dst.A).Should(Equal(123))
	Ω(dst.B).Should(Equal("str"))
	Ω(dst.C).Should(Equal(12.2))
}

func TestEnvDecoderNested(t *testing.T) {
	RegisterTestingT(t)

	type B struct {
		N int
	}
	type C struct {
		S string
		N int
	}
	type D struct {
		B
		N        int
		Nested   C
		Embedded struct {
			N int
		}
	}

	os.Setenv("N", "1")
	os.Setenv("NESTED_N", "30")
	os.Setenv("EMBEDDED_N", "1984")

	dst := D{}
	d := NewEnvDecoder("")
	d.Decode(&dst)

	Ω(dst.N).Should(Equal(1))
	Ω(dst.B.N).Should(Equal(0)) // B.N is shadowed by N
	Ω(dst.Nested.N).Should(Equal(30))
	Ω(dst.Embedded.N).Should(Equal(1984))
	Ω(dst.Nested.S).Should(Equal(""))
}
