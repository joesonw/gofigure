package realworld_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/joesonw/gofigure"
	"github.com/joesonw/gofigure/feature"
)

var _ = Describe("All Realworld Tests", func() {
	It("should resolve after load (can reference  another file)", func() {
		loader := gofigure.New().WithFeatures(
			feature.Template(),
			feature.Reference(),
		)
		Expect(loader.Load("app.yaml", []byte(`env: dev
port: 8080
host: localhost
listen: !tpl |
  {{ config "app.host" }}:{{ config "app.port" }}
name: !ref app.product
`))).To(BeNil())
		Expect(loader.Load("app.yaml", []byte(`env: prod 
product: test
port: 80`))).To(BeNil())
		var listen, name string
		Expect(loader.Get(context.Background(), "app.listen", &listen)).To(BeNil())
		Expect(loader.Get(context.Background(), "app.name", &name)).To(BeNil())
		Expect(listen).To(Equal("localhost:80"))
		Expect(name).To(Equal("test"))
	})
})
