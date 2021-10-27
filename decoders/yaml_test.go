package decoders

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/gomega"
)

func TestNewYamlFileDecoder(t *testing.T) {
	t.Run("interface", func(t *testing.T) {
		RegisterTestingT(t)

		var v interface{}
		d := NewYamlFileDecoder("./config_test.yaml")
		Ω(d.Decode(&v)).Should(BeNil())

		expV := map[interface{}]interface{}{
			"simple": map[interface{}]interface{}{
				"num":  1,
				"str":  "test text",
				"arr":  []interface{}{"a", "b", "d"},
				"bool": true,
			},
			"nested": map[interface{}]interface{}{
				"one": map[interface{}]interface{}{
					"two": []interface{}{
						map[interface{}]interface{}{
							"three":  3,
							"common": "c-3",
						},
						map[interface{}]interface{}{
							"four":   4,
							"common": "c-4",
						},
					},
					"five": []interface{}{6, 7, 8},
				},
				"done": "we are done",
			},
			"names": map[interface{}]interface{}{
				"lowercase":   "just a lowercase key",
				"Capitalized": "a Capitalized key",
				"camelCase":   "a camelCase key",
				"PascalCase":  "PascalCase is cool",
				"snake_case":  "try snake_case",
				"kebab-case":  "kebab-case is always a pleasure to look at",
			},
		}

		Ω(cmp.Diff(v, expV)).Should(BeEmpty())
		Ω(v).Should(Equal(expV))
	})

	t.Run("struct", func(t *testing.T) {
		RegisterTestingT(t)

		var v testYamlConfig
		d := NewYamlFileDecoder("./config_test.yaml")
		Ω(d.Decode(&v)).Should(BeNil())

		expV := testYamlConfig{
			Simple: testYamlConfig_Simple{
				Bool: true,
				Num:  1,
				Str:  "test text",
				Arr:  []string{"a", "b", "d"},
			},
			Nested: testYamlConfig_Nested{
				One: testYamlConfig_Nested_One{
					Two: []testYamlConfig_Nested_One_Two{
						{Three: 3, Common: "c-3"},
						{Four: 4, Common: "c-4"},
					},
					Five: []int{6, 7, 8},
				},
				Done: "we are done",
			},
			Names: testYamlConfig_Names{
				Lowercase:   "just a lowercase key",
				Capitalized: "a Capitalized key",
				CamelCase:   "a camelCase key",
				PascalCase:  "PascalCase is cool",
				SnakeCase:   "try snake_case",
				KebabCase:   "kebab-case is always a pleasure to look at",
			},
		}

		Ω(cmp.Diff(v, expV)).Should(BeEmpty())
		Ω(v).Should(Equal(expV))
	})
}

type testYamlConfig struct {
	Simple testYamlConfig_Simple
	Nested testYamlConfig_Nested
	Names  testYamlConfig_Names
}

type testYamlConfig_Simple struct {
	Bool bool
	Num  int
	Str  string
	Arr  []string
}

type testYamlConfig_Nested struct {
	One  testYamlConfig_Nested_One
	Done string
}

type testYamlConfig_Nested_One struct {
	Two  []testYamlConfig_Nested_One_Two
	Five []int
}

type testYamlConfig_Nested_One_Two struct {
	Three  int
	Four   int
	Common string
}

type testYamlConfig_Names struct {
	Lowercase   string
	Capitalized string `yaml:"Capitalized"`
	CamelCase   string `yaml:"camelCase"`
	PascalCase  string `yaml:"PascalCase"`
	SnakeCase   string `yaml:"snake_case"`
	KebabCase   string `yaml:"kebab-case"`
}
