package decoders

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestStructDecodeKey(t *testing.T) {
	RegisterTestingT(t)

	src := struct {
		A struct {
			B struct {
				N int
				S string
			}
			N int
		}
		D struct {
			N int
		}
		S string
	}{}
	src.A.B.N = 1
	src.A.B.S = "A.B.S val"
	src.A.N = 2
	src.D.N = 3
	src.S = "root"

	d := NewStructDecoder(&src, "", "").(*kvwrapper).store

	dStr := ""
	dInt := 0

	Ω(d.DecodeKey("S", &dStr)).ShouldNot(HaveOccurred())
	Ω(dStr).Should(Equal("root"))

	Ω(d.DecodeKey("A.B.S", &dStr)).ShouldNot(HaveOccurred())
	Ω(dStr).Should(Equal("A.B.S val"))

	Ω(d.DecodeKey("A.B.N", &dInt)).ShouldNot(HaveOccurred())
	Ω(dInt).Should(Equal(1))

	Ω(d.DecodeKey("A.N", &dInt)).ShouldNot(HaveOccurred())
	Ω(dInt).Should(Equal(2))

	Ω(d.DecodeKey("D.N", &dInt)).ShouldNot(HaveOccurred())
	Ω(dInt).Should(Equal(3))
}

func TestStructDecoderBasic(t *testing.T) {
	RegisterTestingT(t)

	src := struct {
		A int     `ttt:"aaa"`
		B string  `ttt:"bbb"`
		C float64 `ttt:"ccc"`
	}{A: 231, B: "asd", C: 1.2}

	d1 := NewStructDecoder(&src, "", "")
	dst1 := struct{ A int }{}
	d1.Decode(&dst1)
	Ω(dst1.A).Should(Equal(231))

	d2 := NewStructDecoder(&src, "ttt", "json")
	dst2 := struct {
		F1 string `json:"bbb"`
		F2 int    `json:"aaa"`
	}{}
	d2.Decode(&dst2)
	Ω(dst2.F1).Should(Equal("asd"))
	Ω(dst2.F2).Should(Equal(231))
}

func TestStructDecoderNested(t *testing.T) {
	RegisterTestingT(t)

	type B struct {
		N int `json:"n"`
	}
	type D struct {
		B
		N        int
		Embedded struct {
			N int
		} `json:"e"`
	}

	src := D{B: B{N: 1}, N: 2, Embedded: struct{ N int }{N: 3}}
	d1 := NewStructDecoder(&src, "json", "ttt")
	dst1 := struct {
		F1 int `ttt:"n"`
		F2 int `ttt:"e.N"`
		F3 int `ttt:"N"`
	}{}
	d1.Decode(&dst1)
	Ω(dst1.F1).Should(Equal(1))
	Ω(dst1.F2).Should(Equal(3))
	Ω(dst1.F3).Should(Equal(2))

	d2 := NewStructDecoder(&src, "json", "json")
	dst2 := struct {
		B
		E struct {
			N int `json:"XXX"`
		}
	}{}
	d2.Decode(&dst2)
	Ω(dst2.B.N).Should(Equal(1))
	Ω(dst2.E.N).Should(Equal(0))

	dst3 := struct {
		B
		E struct {
			N int `json:"XXX"`
		}
		F1 int `json:"n"`
	}{}
	dst3.E.N = 333
	d2.Decode(&dst3)
	Ω(dst3.F1).Should(Equal(1))
	Ω(dst3.B.N).Should(Equal(0)) // `B.N` is shadowed by F1
	Ω(dst3.E.N).Should(Equal(333))
}
