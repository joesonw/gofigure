package gofigure

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Node", func() {
	It("create node with options", func() {
		parent := &Node{}
		node := createNodeWithOptions(
			NodeFilepath("filepath"),
			NodeParent(parent),
			NodeSequenceIndex(123),
			NodeMappingKey("key"),
		)
		Expect(node.Filepath()).To(Equal("filepath"))
		Expect(node.parent).To(Equal(parent))
		Expect(node.sequenceIndex).To(Equal(123))
		Expect(node.mappingKey).To(Equal("key"))
		Expect(node.hasSequenceIndex).To(BeTrue())
		Expect(node.hasMappingKey).To(BeTrue())
	})

	It("should create scalar node", func() {
		node := NewScalarNode("value")
		Expect(node.kind).To(Equal(yaml.ScalarNode))
		Expect(node.value).To(Equal("value"))
	})

	It("should create mapping node", func() {
		value := NewScalarNode("value")
		node := NewMappingNode(map[string]*Node{
			"key": value,
		})
		Expect(node.kind).To(Equal(yaml.MappingNode))
		Expect(node.mappingNodes).To(And(
			HaveLen(1),
			HaveKeyWithValue("key", Equal(value)),
		))
	})

	It("should create sequence node", func() {
		value := NewScalarNode("value")
		node := NewSequenceNode([]*Node{
			value,
		})
		Expect(node.kind).To(Equal(yaml.SequenceNode))
		Expect(node.sequenceNodes).To(And(
			HaveLen(1),
			ContainElement(Equal(value)),
		))
	})

	It("should get filepath", func() {
		node := &Node{
			parent: &Node{
				filepath: "filepath",
			},
		}
		Expect(node.Filepath()).To(Equal("filepath"))

		node = &Node{
			parent: &Node{},
		}
		Expect(node.Filepath()).To(Equal(""))
	})

	It("should get key path", func() {
		node := &Node{
			parent: &Node{
				hasMappingKey: true,
				mappingKey:    "some",
				parent: &Node{
					hasMappingKey: true,
					mappingKey:    "parent",
					parent:        &Node{},
				},
			},
			hasSequenceIndex: true,
			sequenceIndex:    2,
		}
		Expect(node.Keypath()).To(Equal("parent.some[2]"))
	})

	It("should get mapping child", func() {
		node := &Node{
			kind: yaml.MappingNode,
			mappingNodes: map[string]*Node{
				"key": {},
			},
		}

		value, err := node.GetMappingChild("key")
		Expect(err).To(BeNil())
		Expect(value).To(Equal(node.mappingNodes["key"]))

		value, err = node.GetMappingChild("not-exist")
		Expect(err).To(BeNil())
		Expect(value).To(BeNil())
	})

	It("should not get mapping child", func() {
		node := &Node{
			mappingNodes: map[string]*Node{},
		}

		_, err := node.GetMappingChild("key")
		Expect(err).To(MatchError("\"\" is not a mapping node"))
	})

	It("should get sequence child", func() {
		node := &Node{
			kind: yaml.SequenceNode,
			sequenceNodes: []*Node{
				{},
			},
		}

		value, err := node.GetSequenceChild(0)
		Expect(err).To(BeNil())
		Expect(value).To(Equal(node.sequenceNodes[0]))

		value, err = node.GetSequenceChild(1)
		Expect(err).To(BeNil())
		Expect(value).To(BeNil())
	})

	It("should not get sequence child", func() {
		node := &Node{
			sequenceNodes: []*Node{},
		}

		_, err := node.GetSequenceChild(0)
		Expect(err).To(MatchError("\"\" is not a sequence node"))
	})

	It("should get deep value", func() {
		value := &Node{}
		node := &Node{
			kind: yaml.MappingNode,
			mappingNodes: map[string]*Node{
				"nested": {
					kind: yaml.MappingNode,
					mappingNodes: map[string]*Node{
						"key": {
							kind:          yaml.SequenceNode,
							sequenceNodes: []*Node{value},
						},
					},
				},
			},
		}

		result, err := node.GetDeep("nested.key[0]")
		Expect(err).To(BeNil())
		Expect(result).To(Equal(value))

		result, err = node.GetDeep("nested.not.exist")
		Expect(err).To(BeNil())
		Expect(result).To(BeNil())

		result, err = node.GetDeep("")
		Expect(err).To(BeNil())
		Expect(result).To(Equal(node))

		_, err = node.GetDeep("a.")
		Expect(err).To(MatchError("unable to parse path \"a.\": a.: invalid path"))
	})

	It("should to yaml node", func() {
		node := &Node{
			kind: yaml.MappingNode,
			mappingNodes: map[string]*Node{
				"key": {
					kind: yaml.SequenceNode,
					sequenceNodes: []*Node{
						nil,
						{
							kind:  yaml.ScalarNode,
							value: "value",
						},
					},
				},
				"not": nil,
			},
		}

		Expect(node.ToYAMLNode()).To(PointTo(MatchFields(IgnoreExtras, Fields{
			"Kind": Equal(yaml.MappingNode),
			"Content": And(
				HaveLen(2),
				ConsistOf(
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Kind":  Equal(yaml.ScalarNode),
						"Value": Equal("key"),
					})),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Kind": Equal(yaml.SequenceNode),
						"Content": And(
							HaveLen(1),
							ConsistOf(
								PointTo(MatchFields(IgnoreExtras, Fields{
									"Kind":  Equal(yaml.ScalarNode),
									"Value": Equal("value"),
								})),
							),
						),
					})),
				),
			),
		})))
	})
})

var _ = DescribeTable("Node BoolValue", func(value string, b bool) {
	node := &Node{
		kind:  yaml.ScalarNode,
		value: "true",
	}
	result, err := node.BoolValue()
	Expect(result).To(BeTrue())
	Expect(err).To(BeNil())
},
	boolValuePathEntry("true", true),
	boolValuePathEntry("True", true),
	boolValuePathEntry("TRUE", true),
	boolValuePathEntry("false", false),
	boolValuePathEntry("False", false),
	boolValuePathEntry("FALSE", false),
	boolValuePathEntry("1", true),
	boolValuePathEntry("0", false),
	boolValuePathEntry("y", true),
	boolValuePathEntry("n", false),
	boolValuePathEntry("t", true),
	boolValuePathEntry("f", false),
	boolValuePathEntry("Y", true),
	boolValuePathEntry("N", false),
	boolValuePathEntry("T", true),
	boolValuePathEntry("F", false),
	boolValuePathEntry("yes", true),
	boolValuePathEntry("no", false),
	boolValuePathEntry("YES", true),
	boolValuePathEntry("NO", false),
	boolValuePathEntry("Yes", true),
	boolValuePathEntry("No", false),
	boolValuePathEntry("on", true),
	boolValuePathEntry("off", false),
	boolValuePathEntry("ON", true),
	boolValuePathEntry("OFF", false),
	boolValuePathEntry("On", true),
	boolValuePathEntry("Off", false),
)

func boolValuePathEntry(value string, b bool) TableEntry {
	return Entry(fmt.Sprintf("bool value %q should be %v", value, b), value, b)
}
