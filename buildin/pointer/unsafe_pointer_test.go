package pointer

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestUnsafePointer(t *testing.T) {
	// Create an integer array
	arr := [5]int{10, 20, 30, 40, 50}

	// Get the pointer to the first element of the array
	p := unsafe.Pointer(&arr[0])

	// Print the original array and the value of the first element
	fmt.Println("Original array:", arr)       // [10 20 30 40 50]
	fmt.Println("First element:", *(*int)(p)) // 10

	// Calculate the pointer to the next element
	sizeOfInt := unsafe.Sizeof(arr[0])
	p = unsafe.Pointer(uintptr(p) + sizeOfInt)

	// Print the value of the second element
	fmt.Println("Second element:", *(*int)(p)) // 20

	// Continue calculating and printing the values of subsequent elements

	// To demonstrate garbage collection issues, create a pointer and set it to nil
	ptr := &arr[0]
	fmt.Println("Pointer before:", *ptr) // 10
	ptr = nil
	fmt.Println("Pointer after nil:", ptr) // <nil>

	// Even though the pointer is set to nil, the array still exists because it is within the scope of the main function
	fmt.Println("Array after nil pointer:", arr) // [10 20 30 40 50]
}
