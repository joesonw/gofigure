package gofigure

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Loader", func() {
	It("should Load", func() {
		loader := New()
		Expect(loader.Load("config/app.yaml", []byte(`env: dev
name: John`))).To(BeNil())
		Expect(loader.root.kind).To(Equal(yaml.MappingNode))
		Expect(loader.root.mappingNodes["config"].mappingNodes["app"].mappingNodes["env"].value).To(Equal("dev"))
		Expect(loader.root.mappingNodes["config"].mappingNodes["app"].mappingNodes["name"].value).To(Equal("John"))
		Expect(loader.loadedFiles["config/app"]).To(BeTrue())

		Expect(loader.Load("config/app.yaml", []byte(`env: test`))).To(BeNil())
		Expect(loader.root.mappingNodes["config"].mappingNodes["app"].mappingNodes["env"].value).To(Equal("test"))
		Expect(loader.root.mappingNodes["config"].mappingNodes["app"].mappingNodes["name"].value).To(Equal("John"))

		Expect(loader.Load("another.yaml", []byte(`value: another`))).To(BeNil())
		Expect(loader.root.kind).To(Equal(yaml.MappingNode))
		Expect(loader.root.mappingNodes["another"].mappingNodes["value"].value).To(Equal("another"))
		Expect(loader.loadedFiles["another"]).To(BeTrue())
	})

	It("should Get", func() {
		loader := New()
		Expect(loader.Load("config/app.yaml", []byte(`array:
- 213
- 456`))).To(BeNil())
		var value int
		Expect(loader.Get(context.Background(), "config.app.array[0]", &value)).To(BeNil())
		Expect(value).To(Equal(213))
		Expect(loader.Get(context.Background(), "config.app.array[1]", &value)).To(BeNil())
		Expect(value).To(Equal(456))
		value = 789
		Expect(loader.Get(context.Background(), "config.app.array[2]", &value)).To(BeNil())
		Expect(value).To(Equal(789))
	})

})
