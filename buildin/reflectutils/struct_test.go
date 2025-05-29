package reflectutils

import (
	"fmt"
	"reflect"
	"testing"
)

type Manager struct {
	User
	title string
}

func TestHasExportedFields(t *testing.T) {
	managerType := reflect.TypeOf(Manager{})
	hasExported := HasExportedFields(managerType)
	fmt.Println(hasExported) // 输出: true
}
