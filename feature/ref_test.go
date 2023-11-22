package feature_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/joesonw/gofigure"
	"github.com/joesonw/gofigure/feature"
)

var _ = Describe("!ref", func() {
	It("should work", func() {
		loader := gofigure.New().WithFeatures(
			feature.Reference(),
		)
		Expect(loader.Load("app.yaml", []byte(`complex:
  first: John
  last: Doe
test_value: !ref app.complex`))).To(BeNil())
		var first string
		var last string
		Expect(loader.Get(context.Background(), "app.test_value.first", &first)).To(BeNil())
		Expect(loader.Get(context.Background(), "app.test_value.last", &last)).To(BeNil())
		Expect(first).To(Equal("John"))
		Expect(last).To(Equal("Doe"))
	})
})
