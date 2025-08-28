package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetTag(t *testing.T) {
	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		email string `json:"email"`
	}

	userType := reflect.TypeOf(User{})
	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fmt.Printf("Field: %s, Tag: %s\n", field.Name, field.Tag.Get("json"))
	}

	// Field: ID, Tag: id
	// Field: Name, Tag: name
	// Field: email, Tag: email
}
