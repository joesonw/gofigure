package feature_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psanford/memfs"

	"github.com/joesonw/gofigure"
	"github.com/joesonw/gofigure/feature"
)

var _ = Describe("!include", func() {
	It("should work", func() {
		fs := memfs.New()
		source := `name:
  first: John
  last: Doe`
		Expect(fs.WriteFile("external.yaml", []byte(source), 0644)).To(BeNil())
		loader := gofigure.New().WithFeatures(
			feature.Include(fs),
			feature.Template(),
		)
		Expect(loader.Load("app.yaml", []byte(`file: external
test_value_plain: !include
  file: 
    path: !tpl |
      {{ config "app.file" }}.yaml
test_value_complex: !include
  file:
    path: external.yaml
    parse: true
test_value_simple:
  first: !include
    file:
      path: external.yaml
      parse: true
      key: name.first`))).To(BeNil())
		var plainValue string
		var complexValue struct {
			Name struct {
				First string `yaml:"first"`
				Last  string `yaml:"last"`
			} `yaml:"name"`
		}
		var simpleValue string
		Expect(loader.Get(context.Background(), "app.test_value_plain", &plainValue)).To(BeNil())
		Expect(loader.Get(context.Background(), "app.test_value_complex", &complexValue)).To(BeNil())
		Expect(loader.Get(context.Background(), "app.test_value_simple.first", &simpleValue)).To(BeNil())
		Expect(plainValue).To(Equal(source))
		Expect(complexValue.Name.First).To(Equal("John"))
		Expect(complexValue.Name.Last).To(Equal("Doe"))
		Expect(simpleValue).To(Equal("John"))
	})

	It("should work when files circular referencing", func() {
		fs := memfs.New()
		source := `first: !ref app.name.first`
		Expect(fs.WriteFile("external.yaml", []byte(source), 0644)).To(BeNil())
		loader := gofigure.New().WithFeatures(
			feature.Include(fs),
			feature.Reference(),
		)
		Expect(loader.Load("app.yaml", []byte(`name:
  first: John
value: !include
  file:
    path: external.yaml
    parse: true
    key: first`))).To(BeNil())
		var value string
		Expect(loader.Get(context.Background(), "app.value", &value)).To(BeNil())
		Expect(value).To(Equal("John"))
	})
})
