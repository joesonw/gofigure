package gofigure

type nodeOptions struct {
	filepath         string
	mappingKey       string
	hasMappingKey    bool
	sequenceIndex    int
	hasSequenceIndex bool
	parent           *Node
}

type NodeOption interface {
	apply(*nodeOptions)
}

type nodeOptionFunc func(*nodeOptions)

func (f nodeOptionFunc) apply(o *nodeOptions) {
	f(o)
}

func NodeFilepath(path string) NodeOption {
	return nodeOptionFunc(func(o *nodeOptions) {
		o.filepath = path
	})
}

func NodeSequenceIndex(index int) NodeOption {
	return nodeOptionFunc(func(o *nodeOptions) {
		o.sequenceIndex = index
		o.hasSequenceIndex = true
	})
}

func NodeMappingKey(key string) NodeOption {
	return nodeOptionFunc(func(o *nodeOptions) {
		o.mappingKey = key
		o.hasMappingKey = true
	})
}

func NodeParent(parent *Node) NodeOption {
	return nodeOptionFunc(func(o *nodeOptions) {
		o.parent = parent
	})
}
