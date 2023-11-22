package feature_test

import (
	"context"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/joesonw/gofigure"
	"github.com/joesonw/gofigure/feature"
)

var _ = Describe("!tpl", func() {
	It("should work", func() {
		loader := gofigure.New().WithFeatures(
			feature.Template().Funcs(template.FuncMap{
				"say": func() string {
					return "not today"
				},
			}).Valeus(map[string]any{
				"name": "John",
			}),
		)
		Expect(loader.Load("app.yaml", []byte(`
env: dev
test_value: !tpl This is {{ .name }} says {{ say }} in {{ config "app.env" }}`))).To(BeNil())
		var value string
		Expect(loader.Get(context.Background(), "app.test_value", &value)).To(BeNil())
		Expect(value).To(Equal("This is John says not today in dev"))
	})
})
