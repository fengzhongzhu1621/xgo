package randutils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMD5Hash(t *testing.T) {
	text := "xxx"
	expected := "f561aaf6ef0bf14d4208bb46a4ccb3ad"
	result := GetMD5Hash(text)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetBytesMD5Hash(t *testing.T) {
	text := []byte("xxx")
	expected := "f561aaf6ef0bf14d4208bb46a4ccb3ad"
	result := GetBytesMD5Hash(text)
	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestGetSliceMD5Hash(t *testing.T) {
	strSlice := []string{"xxx", "yyy"}
	expectedStr := "5f221cf63a70ca156f0fe1058e7f250b"
	resultStr, err := GetSliceMD5Hash(strSlice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resultStr != expectedStr {
		t.Errorf("Expected %v, but got %v", expectedStr, resultStr)
	}
	intSlice := []int64{1, 2, 3}
	expectedInt := "55b84a9d317184fe61224bfb4a060fb0"
	resultInt, err := GetSliceMD5Hash(intSlice)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resultInt != expectedInt {
		t.Errorf("Expected %v, but got %v", expectedInt, resultInt)
	}
	invalidSlice := []float64{1.1, 2.2, 3.3}
	_, err = GetSliceMD5Hash(invalidSlice)
	assert.Error(t, err)
	assert.Equal(t, errors.New("illegal type").Error(), err.Error())
}
