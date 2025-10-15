package multierror

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-multierror"
)

func TestAppend(t *testing.T) {
	var result error

	// 模拟多个可能失败的操作
	if err := step1(); err != nil {
		result = multierror.Append(result, err)
	}
	if err := step2(); err != nil {
		result = multierror.Append(result, err)
	}

	// 返回聚合后的错误，如果没有错误则返回 nil
	if result != nil {
		fmt.Println(result.Error())
	}

	// 2 errors occurred:
	// * 第一步失败
	// * 第二步失败

}

func step1() error {
	return errors.New("第一步失败")
}

func step2() error {
	return errors.New("第二步失败")
}
