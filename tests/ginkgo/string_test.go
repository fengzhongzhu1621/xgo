package ginkgo

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("String", func() {

	Describe("TruncateString", func() {
		var s = "helloworld"

		DescribeTable("TruncateString cases", func(expected string, truncatedSize int) {
			assert.Equal(GinkgoT(), expected, TruncateString(s, truncatedSize))
		},
			Entry("truncated size less than real size", "he", 2),
			Entry("truncated size equals to real size", s, 10),
			Entry("truncated size greater than real size", s, 20),
		)
	})
})
