package gofigure

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dotpath", func() {
	It("should parse path a.b[0].c[1].d", func() {
		paths, err := ParseDotPath("a.b[0].c[1].d")
		Expect(err).To(BeNil())
		Expect(paths).To(HaveLen(6))
		Expect(paths[0].Key).To(Equal("a"))
		Expect(paths[1].Key).To(Equal("b"))
		Expect(paths[2].Index).To(Equal(0))
		Expect(paths[3].Key).To(Equal("c"))
		Expect(paths[4].Index).To(Equal(1))
		Expect(paths[5].Key).To(Equal("d"))
	})

	It("should return nil if empty", func() {
		paths, err := ParseDotPath("")
		Expect(err).To(BeNil())
		Expect(paths).To(BeNil())
	})
})

var _ = DescribeTable("Dotpath failed scenarios", func(path string) {
	_, err := ParseDotPath(path)
	Expect(err).To(MatchError(path + ": invalid path"))
},
	dotPathEntry(".b"),
	dotPathEntry("[.a"),
	dotPathEntry("a[b]"),
	dotPathEntry("a..b"),
	dotPathEntry("[["),
	dotPathEntry("]"),
	dotPathEntry("[]"),
	dotPathEntry("a["),
	dotPathEntry("a."),
)

func dotPathEntry(path string) TableEntry {
	return Entry(fmt.Sprintf("path %q", path), path)
}
