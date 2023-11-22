package gofigure

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Node struct {
	kind        yaml.Kind
	style       yaml.Style
	value       string
	tag         string
	anchor      string
	headComment string
	lineComment string
	footComment string
	line        int
	column      int

	filepath string

	parent           *Node
	sequenceIndex    int
	hasSequenceIndex bool
	mappingKey       string
	hasMappingKey    bool

	resolved     bool
	resolvedNode *Node

	mappingNodes  map[string]*Node
	sequenceNodes []*Node
}

func createNodeWithOptions(options ...NodeOption) *Node {
	o := &nodeOptions{}
	for i := range options {
		options[i].apply(o)
	}
	n := &Node{
		filepath:         o.filepath,
		parent:           o.parent,
		sequenceIndex:    o.sequenceIndex,
		hasSequenceIndex: o.hasSequenceIndex,
		mappingKey:       o.mappingKey,
		hasMappingKey:    o.hasMappingKey,
	}
	return n
}

func NewScalarNode(value string, options ...NodeOption) *Node {
	n := createNodeWithOptions(options...)
	n.kind = yaml.ScalarNode
	n.value = value
	return n
}

func NewMappingNode(m map[string]*Node, options ...NodeOption) *Node {
	n := createNodeWithOptions(options...)
	n.kind = yaml.MappingNode
	n.mappingNodes = m
	return n
}

func NewSequenceNode(values []*Node, options ...NodeOption) *Node {
	n := createNodeWithOptions(options...)
	n.kind = yaml.SequenceNode
	n.sequenceNodes = values
	return n
}

func NewNode(node *yaml.Node, options ...NodeOption) *Node {
	n := createNodeWithOptions(options...)
	n.unmarshalNode(node)
	return n
}

func (n *Node) Filepath() string {
	if n.filepath != "" {
		return n.filepath
	}

	if n.parent != nil {
		return n.parent.Filepath()
	}

	return ""
}

func (n *Node) Keypath() string {
	var key string
	cur := n
	for cur != nil {
		if cur.hasSequenceIndex {
			key = fmt.Sprintf("[%d]%s", cur.sequenceIndex, key)
		} else if cur.hasMappingKey {
			dot := "."
			if strings.HasPrefix(key, "[") {
				dot = ""
			}
			key = fmt.Sprintf("%s%s%s", cur.mappingKey, dot, key)
		} else {
			break
		}
		cur = cur.parent
	}
	return strings.TrimSuffix(key, ".")
}

func (n *Node) GetMappingChild(key string) (*Node, error) {
	if n.kind != yaml.MappingNode {
		return nil, fmt.Errorf("%q is not a mapping node", n.Keypath())
	}

	return n.mappingNodes[key], nil
}

func (n *Node) GetSequenceChild(index int) (*Node, error) {
	if n.kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%q is not a sequence node", n.Keypath())
	}
	if index >= len(n.sequenceNodes) {
		return nil, nil
	}
	return n.sequenceNodes[index], nil
}

func (n *Node) GetDeep(path string) (*Node, error) {
	if len(path) == 0 {
		return n, nil
	}

	paths, err := parseDotPath(path)
	if err != nil {
		return nil, fmt.Errorf("unable to parse path %q: %w", path, err)
	}

	current := n
	for _, p := range paths {
		if current == nil {
			return nil, nil
		}
		if p.key != "" { // map
			current, err = current.GetMappingChild(p.key)
			if err != nil {
				return nil, err
			}
		} else { // slice
			current, err = current.GetSequenceChild(p.index)
			if err != nil {
				return nil, err
			}
		}
	}

	return current, nil
}

func (n *Node) MarshalYAML() (interface{}, error) {
	return n.ToYAMLNode(), nil
}

func (n *Node) UnmarshalYAML(value *yaml.Node) error {
	n.unmarshalNode(value)
	return nil
}

func (n *Node) unmarshalNode(node *yaml.Node) {
	n.kind = node.Kind
	n.style = node.Style
	n.value = node.Value
	n.tag = node.Tag
	n.anchor = node.Anchor
	n.headComment = node.HeadComment
	n.lineComment = node.LineComment
	n.footComment = node.FootComment
	n.line = node.Line
	n.column = node.Column
	setNodeValueFromYAML(n, node)
}

func setNodeValueFromYAML(n *Node, node *yaml.Node) {
	switch node.Kind {
	case yaml.DocumentNode, yaml.ScalarNode:
	case yaml.MappingNode:
		n.mappingNodes = make(map[string]*Node)
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			childNode := NewNode(valueNode)
			childNode.parent = n
			childNode.mappingKey = keyNode.Value
			childNode.hasMappingKey = true
			n.mappingNodes[keyNode.Value] = childNode
		}
	case yaml.SequenceNode:
		n.sequenceNodes = make([]*Node, len(node.Content))
		for i := range node.Content {
			childNode := NewNode(node.Content[i])
			childNode.parent = n
			childNode.sequenceIndex = i
			childNode.hasSequenceIndex = true
			n.sequenceNodes[i] = childNode
		}
	case yaml.AliasNode:
		n.kind = node.Alias.Kind
		n.value = node.Alias.Value
		setNodeValueFromYAML(n, node.Alias)
	}
}

func (n *Node) ToYAMLNode() *yaml.Node {
	var node yaml.Node
	node.Kind = n.kind
	node.Style = n.style
	node.Value = n.value
	node.Tag = n.tag
	node.Anchor = n.anchor
	node.HeadComment = n.headComment
	node.LineComment = n.lineComment
	node.FootComment = n.footComment
	node.Line = n.line
	node.Column = n.column

	switch n.kind {
	case yaml.DocumentNode, yaml.AliasNode, yaml.ScalarNode:
	case yaml.MappingNode:
		node.Content = nil
		for key, childNode := range n.mappingNodes {
			if childNode == nil {
				continue
			}
			keyNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: key,
			}
			node.Content = append(node.Content, keyNode, childNode.ToYAMLNode())
		}
	case yaml.SequenceNode:
		node.Content = nil
		for _, childNode := range n.sequenceNodes {
			if childNode == nil {
				continue
			}
			node.Content = append(node.Content, childNode.ToYAMLNode())
		}
	}

	return &node
}

func (n *Node) Kind() yaml.Kind {
	return n.kind
}

func (n *Node) Style() yaml.Style {
	return n.style
}

func (n *Node) Value() string {
	return n.value
}

func (n *Node) Tag() string {
	return n.tag
}

func (n *Node) Anchor() string {
	return n.anchor
}

func (n *Node) HeadComment() string {
	return n.headComment
}

func (n *Node) LineComment() string {
	return n.lineComment
}

func (n *Node) FootComment() string {
	return n.footComment
}

func (n *Node) Line() int {
	return n.line
}

func (n *Node) Column() int {
	return n.column
}

func (n *Node) BoolValue() (bool, error) {
	if n.kind != yaml.ScalarNode {
		return false, fmt.Errorf("%q is not a scalar node", n.Keypath())
	}

	var b bool
	if err := yaml.Unmarshal([]byte(n.value), &b); err != nil {
		return false, err
	}
	return b, nil
}
