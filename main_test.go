package xgo

import (
	"testing"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	// 用于检测测试过程中是否有goroutine泄漏。
	// 会在所有测试用例运行完毕后，检查是否有未退出的goroutine。如果有，它会打印出泄漏的goroutine信息，并使测试失败。
	//
	// 使用场景
	// 并发测试: 当你在编写涉及并发操作的测试时，可能会不小心留下一些未退出的goroutine，这会导致资源泄漏。使用goleak可以帮助你及时发现这些问题。
	// 长时间运行的测试: 对于一些长时间运行的测试，goroutine泄漏可能会导致内存占用不断增加，最终影响系统的稳定性。goleak可以帮助你监控和检测这些问题。
	goleak.VerifyTestMain(m)
}
