package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	email string `json:"email"`
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.email, // 主动赋值
	})
}

// 导出结构体的私有字段
func TestMarshalPrivateField(t *testing.T) {
	u := User{ID: 1, Name: "Tom", email: "tom@example.com"}
	data, _ := json.Marshal(u)
	fmt.Println(string(data)) // {"id":1,"name":"Tom","email":"tom@example.com"}
}
