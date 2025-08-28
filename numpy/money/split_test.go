package money

import (
	"fmt"
	"log"
	"testing"

	"github.com/Rhymond/go-money"
)

func TestSplit(t *testing.T) {
	// 初始化金额为100分（即1人民币）
	pound := money.New(100, money.CNY) // 1 人民币

	// 分割金额为3份
	parties, err := pound.Split(3)
	if err != nil {
		log.Fatalf("分割金额时出错: %v", err)
	}

	// 输出分割后的每一份
	fmt.Printf("分割后的金额:\n")
	for i, part := range parties {
		fmt.Printf("第%d份: %s\n", i+1, part.Display())
	}

	// 分配金额为33, 33, 33分
	partiesAllocated, err := pound.Allocate(33, 33, 33)
	if err != nil {
		log.Fatalf("分配金额时出错: %v", err)
	}

	// 输出分配后的每一份
	fmt.Printf("\n分配后的金额:\n")
	for i, part := range partiesAllocated {
		fmt.Printf("第%d份: %s\n", i+1, part.Display())
	}

	// 分割后的金额:
	// 第1份: 0.34 元
	// 第2份: 0.33 元
	// 第3份: 0.33 元

	// 分配后的金额:
	// 第1份: 0.34 元
	// 第2份: 0.33 元
	// 第3份: 0.33 元

	// 分割为5份
	parties5, err := pound.Split(5)
	if err != nil {
		log.Fatalf("分割金额为5份时出错: %v", err)
	}
	fmt.Printf("\n分割为5份后的金额:\n")
	for i, part := range parties5 {
		fmt.Printf("第%d份: %s\n", i+1, part.Display())
	}

	// 分配为40分、30分、30分
	partiesAllocated2, err := pound.Allocate(40, 30, 30)
	if err != nil {
		log.Fatalf("分配金额为40, 30, 30分时出错: %v", err)
	}
	fmt.Printf("\n分配为40, 30, 30分后的金额:\n")
	for i, part := range partiesAllocated2 {
		fmt.Printf("第%d份: %s\n", i+1, part.Display())
	}

	// 分割为5份后的金额:
	// 第1份: 0.20 元
	// 第2份: 0.20 元
	// 第3份: 0.20 元
	// 第4份: 0.20 元
	// 第5份: 0.20 元

	// 分配为40, 30, 30分后的金额:
	// 第1份: 0.40 元
	// 第2份: 0.30 元
	// 第3份: 0.30 元
}
