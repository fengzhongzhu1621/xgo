package gofakeit

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

// Person 结构体用于存储生成的假数据
type Person struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func generateFakePerson(seed int64) Person {
	// 设置随机种子
	gofakeit.Seed(seed)

	// 生成假数据
	return Person{
		Name:    gofakeit.Name(),
		Email:   gofakeit.Email(),
		Phone:   gofakeit.Phone(),
		Address: gofakeit.Address().Address,
	}
}

func TestBase(t *testing.T) {
	// 生成假数据
	person := generateFakePerson(0)

	// 打印结构化数据
	fmt.Println("Generated Person Data:")
	fmt.Printf("%+v\n", person)

	// 打印JSON格式数据
	jsonData, err := json.MarshalIndent(person, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	fmt.Println("\nJSON Format:")
	fmt.Println(string(jsonData))
}
