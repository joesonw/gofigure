package gofigure

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Merge", func() {
	It("should not merge document node", func() {
		var nodeA yaml.Node
		Expect(yaml.Unmarshal([]byte(`name: John"`), &nodeA)).To(BeNil())

		var nodeB yaml.Node
		Expect(yaml.Unmarshal([]byte(`age: 123`), &nodeB)).To(BeNil())
		_, err := MergeNodes(NewNode(&nodeA), NewNode(&nodeB))
		Expect(err).To(And(Not(BeNil()), MatchError("document node cannot be merged")))
	})

	It("should not merge unmatched node", func() {
		var nodeA yaml.Node
		Expect(yaml.Unmarshal([]byte(`name: John"`), &nodeA)).To(BeNil())

		var nodeB yaml.Node
		Expect(yaml.Unmarshal([]byte(`name: { age: 1 }`), &nodeB)).To(BeNil())
		_, err := MergeNodes(NewNode(nodeA.Content[0]), NewNode(nodeB.Content[0]))
		Expect(err).To(And(Not(BeNil()), MatchError("cannot merge 8 with 4")))
	})

	It("should merge", func() {
		var nodeA yaml.Node
		Expect(yaml.Unmarshal([]byte(`name: Doe
tagged: value
map:
  key: 123`), &nodeA)).To(BeNil())

		var nodeB yaml.Node
		Expect(yaml.Unmarshal([]byte(`realname: &realname John
name: *realname
tagged: !tagged override
age: 123
map: !replace
  key2: 456`), &nodeB)).To(BeNil())
		result, err := MergeNodes(NewNode(nodeA.Content[0]), NewNode(nodeB.Content[0]))
		Expect(err).To(BeNil())

		b, err := yaml.Marshal(result.ToYAMLNode())
		Expect(err).To(BeNil())

		var s struct {
			Name   string         `yaml:"name"`
			Age    int            `yaml:"age"`
			Tagged string         `yaml:"tagged"`
			Map    map[string]int `yaml:"map"`
		}
		Expect(yaml.Unmarshal(b, &s)).To(BeNil())
		Expect(s.Name).To(Equal("John"))
		Expect(s.Age).To(Equal(123))
		Expect(s.Tagged).To(Equal("override"))
		Expect(s.Map).To(And(
			HaveLen(1),
			HaveKeyWithValue("key2", 456),
		))
	})

	It("should pack node in nested keys", func() {
		b, err := yaml.Marshal(map[string]any{
			"name": "John",
			"birthdate": map[string]any{
				"year":  1970,
				"month": 1,
				"day":   1,
			},
		})
		Expect(err).To(BeNil())

		var node yaml.Node
		Expect(yaml.Unmarshal(b, &node)).To(BeNil())

		result := PackNodeInNestedKeys(NewNode(node.Content[0]), "nested", "keys")
		resultNode := result.ToYAMLNode()

		b, err = yaml.Marshal(resultNode)
		Expect(err).To(BeNil())

		var s struct {
			Nested struct {
				Keys struct {
					Name      string `yaml:"name"`
					Birthdate struct {
						Year  int `yaml:"year"`
						Month int `yaml:"month"`
						Day   int `yaml:"day"`
					}
				} `yaml:"keys"`
			} `yaml:"nested"`
		}
		Expect(yaml.Unmarshal(b, &s)).To(BeNil())
		Expect(s.Nested.Keys.Name).To(Equal("John"))
		Expect(s.Nested.Keys.Birthdate.Year).To(BeEquivalentTo(1970))
		Expect(s.Nested.Keys.Birthdate.Month).To(BeEquivalentTo(1))
		Expect(s.Nested.Keys.Birthdate.Day).To(BeEquivalentTo(1))
	})

	It("should merge two nodes", func() {
		a, err := yaml.Marshal(map[string]any{
			"str": "abc",
			"array": []any{
				1,
				2,
			},
			"map": map[string]any{
				"a": "abc",
				"b": 123,
				"d": []any{
					"a",
				},
			},
		})
		Expect(err).To(BeNil())

		b, err := yaml.Marshal(map[string]any{
			"str": "def",
			"array": []any{
				3,
			},
			"map": map[string]any{
				"b": 456,
				"c": "def",
				"d": []any{
					"b",
				},
			},
		})
		Expect(err).To(BeNil())

		var nodeA, nodeB yaml.Node
		Expect(yaml.Unmarshal(a, &nodeA)).To(BeNil())
		Expect(yaml.Unmarshal(b, &nodeB)).To(BeNil())
		result, err := MergeNodes(NewNode(nodeA.Content[0]), NewNode(nodeB.Content[0]))
		Expect(err).To(BeNil())

		resultNode := result.ToYAMLNode()
		Expect(err).To(BeNil())

		resultBytes, err := yaml.Marshal(resultNode)
		Expect(err).To(BeNil())

		var s struct {
			Str   string `yaml:"str"`
			Array []int  `yaml:"array"`
			Map   struct {
				A string   `yaml:"a"`
				B int      `yaml:"b"`
				C string   `yaml:"c"`
				D []string `yaml:"d"`
			}
		}
		Expect(yaml.Unmarshal(resultBytes, &s)).To(BeNil())
		Expect(s.Str).To(Equal("def"))
		Expect(s.Array).To(Equal([]int{1, 2, 3}))
		Expect(s.Map.A).To(Equal("abc"))
		Expect(s.Map.B).To(BeEquivalentTo(456))
		Expect(s.Map.C).To(Equal("def"))
		Expect(s.Map.D).To(Equal([]string{"a", "b"}))
	})
})
