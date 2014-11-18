package reflectx

import (
	"reflect"
	"strings"
	"testing"
)

func ival(v reflect.Value) int {
	return v.Interface().(int)
}

func TestBasic(t *testing.T) {
	type Foo struct {
		A int
		B int
		C int
	}

	f := Foo{1, 2, 3}
	fv := reflect.ValueOf(f)
	m := NewMapper("")

	v := m.FieldByName(fv, "A")
	if ival(v) != f.A {
		t.Errorf("Expecting %d, got %d", ival(v), f.A)
	}
	v = m.FieldByName(fv, "B")
	if ival(v) != f.B {
		t.Errorf("Expecting %d, got %d", f.B, ival(v))
	}
	v = m.FieldByName(fv, "C")
	if ival(v) != f.C {
		t.Errorf("Expecting %d, got %d", f.C, ival(v))
	}
}

func TestEmbedded(t *testing.T) {
	type Foo struct {
		A int
	}

	type Bar struct {
		Foo
		B int
	}

	type Baz struct {
		A int
		Bar
	}

	m := NewMapper("")

	z := Baz{}
	z.A = 1
	z.B = 2
	z.Bar.Foo.A = 3
	zv := reflect.ValueOf(z)

	v := m.FieldByName(zv, "A")
	if ival(v) != z.A {
		t.Errorf("Expecting %d, got %d", ival(v), z.A)
	}
	v = m.FieldByName(zv, "B")
	if ival(v) != z.B {
		t.Errorf("Expecting %d, got %d", ival(v), z.B)
	}
}

func TestMapping(t *testing.T) {
	type Person struct {
		ID           int
		Name         string
		WearsGlasses bool `db:"wears_glasses"`
	}

	m := NewMapperFunc("db", strings.ToLower)
	p := Person{1, "Jason", true}
	mapping := m.TypeMap(reflect.TypeOf(p))

	for _, key := range []string{"id", "name", "wears_glasses"} {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}
	}

	type SportsPerson struct {
		Weight int
		Age    int
		Person
	}
	s := SportsPerson{Weight: 100, Age: 30, Person: p}
	mapping = m.TypeMap(reflect.TypeOf(s))
	for _, key := range []string{"id", "name", "wears_glasses", "weight", "age"} {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}

	}

	type RugbyPlayer struct {
		Position   int
		IsIntense  bool `db:"is_intense"`
		IsAllBlack bool `db:"-"`
		SportsPerson
	}
	r := RugbyPlayer{12, true, false, s}
	mapping = m.TypeMap(reflect.TypeOf(r))
	for _, key := range []string{"id", "name", "wears_glasses", "weight", "age", "position", "is_intense"} {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}
	}

	if _, ok := mapping["isallblack"]; ok {
		t.Errorf("Expecting to ignore `IsAllBlack` field")
	}

	type EmbeddedLiteral struct {
		Embedded struct {
			Person   string
			Position int
		}
		IsIntense bool
	}

	e := EmbeddedLiteral{}
	mapping = m.TypeMap(reflect.TypeOf(e))
	for _, key := range []string{"person", "position", "isintense"} {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}
	}
}

func TestBasicReducer(t *testing.T) {
	type Foo struct {
		A int
		B int
		C int
	}

	f := Foo{1, 2, 3}
	fv := reflect.ValueOf(f)
	m := NewMapper("")
	m.SetReduceFunc(DelimiterKeyReducer("_"))

	v := m.FieldByName(fv, "A")
	if ival(v) != f.A {
		t.Errorf("Expecting %d, got %d", ival(v), f.A)
	}
	v = m.FieldByName(fv, "B")
	if ival(v) != f.B {
		t.Errorf("Expecting %d, got %d", f.B, ival(v))
	}
	v = m.FieldByName(fv, "C")
	if ival(v) != f.C {
		t.Errorf("Expecting %d, got %d", f.C, ival(v))
	}
}

func TestReducerEmbedded(t *testing.T) {
	type Foo struct {
		A int
	}

	type Bar struct {
		Foo
		B int
	}

	type Baz struct {
		A int
		Bar
	}

	m := NewMapper("")
	m.SetReduceFunc(DelimiterKeyReducer("."))

	z := Baz{}
	z.A = 1
	z.B = 2
	z.Bar.Foo.A = 3
	zv := reflect.ValueOf(z)

	v := m.FieldByName(zv, "A")
	if ival(v) != z.A {
		t.Errorf("Expecting %d, got %d", ival(v), z.A)
	}
	v = m.FieldByName(zv, "B")
	if ival(v) != z.B {
		t.Errorf("Expecting %d, got %d", ival(v), z.B)
	}
}

func TestReducerNested(t *testing.T) {
	type Person struct {
		ID           int
		Name         string
		WearsGlasses bool `db:"W_G"`
	}

	m := NewMapperFunc("db", strings.ToUpper)
	m.SetReduceFunc(DelimiterKeyReducer("_"))

	p := Person{1, "Jason", true}
	mapping := m.TypeMap(reflect.TypeOf(p))

	for _, key := range []string{"ID", "NAME", "W_G"} {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}
	}

	type SportsPerson struct {
		Weight int
		Age    int
		Person
		Physician Person `db:"DOCTOR"`
	}
	s := SportsPerson{Weight: 100, Age: 30, Person: p, Physician: Person{ID: 12}}
	mapping = m.TypeMap(reflect.TypeOf(s))
	keys := []string{"ID", "NAME", "W_G", "WEIGHT", "AGE",
		"DOCTOR_ID", "DOCTOR_NAME", "DOCTOR_W_G"}
	for _, key := range keys {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}

	}

	type RugbyPlayer struct {
		Psychologist Person `db:"DOCTOR"`
		IsIntense    bool   `db:"is_intense"`
		IsAllBlack   bool   `db:"-"`
		SportsPerson
	}
	r := RugbyPlayer{Psychologist: Person{ID: 23, Name: "JJ"},
		SportsPerson: s, IsIntense: true, IsAllBlack: false}
	mapping = m.TypeMap(reflect.TypeOf(r))
	keys = []string{"ID", "NAME", "W_G", "WEIGHT", "AGE",
		"DOCTOR_ID", "DOCTOR_NAME", "DOCTOR_W_G", "is_intense"}
	for _, key := range keys {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}
	}

	if _, ok := mapping["isallblack"]; ok {
		t.Errorf("Expecting to ignore `IsAllBlack` field")
	}

	v := m.FieldByName(reflect.ValueOf(r), "DOCTOR_ID")
	if ival(v) != r.Psychologist.ID {
		t.Errorf("Expecting %d, got %d", r.Psychologist.ID, ival(v))
	}

	type EmbeddedLiteral struct {
		Embedded struct {
			Person   string
			Position int
		}
		IsIntense bool
	}

	e := EmbeddedLiteral{}
	mapping = m.TypeMap(reflect.TypeOf(e))
	keys = []string{"EMBEDDED_PERSON", "EMBEDDED_POSITION", "ISINTENSE"}
	for _, key := range keys {
		if _, ok := mapping[key]; !ok {
			t.Errorf("Expecting to find key %s in mapping but did not.", key)
		}
	}
}

type E1 struct {
	A int
}
type E2 struct {
	E1
	B int
}
type E3 struct {
	E2
	C int
}
type E4 struct {
	E3
	D int
}

func BenchmarkFieldNameL1(b *testing.B) {
	e4 := E4{D: 1}
	for i := 0; i < b.N; i++ {
		v := reflect.ValueOf(e4)
		f := v.FieldByName("D")
		if f.Interface().(int) != 1 {
			b.Fatal("Wrong value.")
		}
	}
}

func BenchmarkFieldNameL4(b *testing.B) {
	e4 := E4{}
	e4.A = 1
	for i := 0; i < b.N; i++ {
		v := reflect.ValueOf(e4)
		f := v.FieldByName("A")
		if f.Interface().(int) != 1 {
			b.Fatal("Wrong value.")
		}
	}
}

func BenchmarkFieldPosL1(b *testing.B) {
	e4 := E4{D: 1}
	for i := 0; i < b.N; i++ {
		v := reflect.ValueOf(e4)
		f := v.Field(1)
		if f.Interface().(int) != 1 {
			b.Fatal("Wrong value.")
		}
	}
}

func BenchmarkFieldPosL4(b *testing.B) {
	e4 := E4{}
	e4.A = 1
	for i := 0; i < b.N; i++ {
		v := reflect.ValueOf(e4)
		f := v.Field(0)
		f = f.Field(0)
		f = f.Field(0)
		f = f.Field(0)
		if f.Interface().(int) != 1 {
			b.Fatal("Wrong value.")
		}
	}
}

func BenchmarkFieldByIndexL4(b *testing.B) {
	e4 := E4{}
	e4.A = 1
	idx := []int{0, 0, 0, 0}
	for i := 0; i < b.N; i++ {
		v := reflect.ValueOf(e4)
		f := FieldByIndexes(v, idx)
		if f.Interface().(int) != 1 {
			b.Fatal("Wrong value.")
		}
	}
}
