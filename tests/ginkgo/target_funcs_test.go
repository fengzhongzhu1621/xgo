package ginkgo

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TargetFuncs", func() {
	BeforeEach(func() {
		println("BeforeEach-2")
	})
	BeforeEach(func() {
		println("BeforeEach-1")
	})
	JustBeforeEach(func() {
		println("JustBeforeEach-1")
	})
	JustBeforeEach(func() {
		println("JustBeforeEach-2")
	})
	JustAfterEach(func() {
		println("JustAfterEach-1")
	})
	JustAfterEach(func() {
		println("JustAfterEach-2")
	})
	AfterEach(func() {
		println("AfterEach-1")
	})
	AfterEach(func() {
		println("AfterEach-2")
	})

	Describe("ReturnInt", func() {
		Context("default", func() {
			var (
				input  int
				result int
			)
			BeforeEach(func() {
				println("BeforeEach in Context")
				input = 1
				result = 1
			})
			AfterEach(func() {
				println("AfterEach in Context")
				input = 0
				result = 0
			})
			It("return value1", func() {
				println("Exec Test1 Case in default content")
				v := ReturnInt(input)
				Expect(v).To(Equal(result))
			})
			It("return value2", func() {
				println("Exec Test2 Case in default content")
				v := ReturnInt(input)
				Expect(v).To(Equal(result))
			})
		})
	})

	Describe("ReturnInt2", func() {
		Context("default2", func() {
			var (
				input  int
				result int
			)
			BeforeEach(func() {
				println("BeforeEach in default2 Context")
				input = 1
				result = 1
			})
			AfterEach(func() {
				println("AfterEach in default2 Context")
				input = 0
				result = 0
			})
			It("return value3", func() {
				println("Exec Test3 Case in default2 content")
				v := ReturnInt(input)
				Expect(v).To(Equal(result))
			})
		})
	})
})
