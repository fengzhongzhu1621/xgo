package network

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestByteEndian(t *testing.T) {
	bytes := []byte{0x78, 0x56, 0x34, 0x12}

	// 0x12345678
	fmt.Printf(
		"LittleEndian: 0x%x\n", binary.LittleEndian.Uint32(bytes),
	)
	// 0x78563412
	fmt.Printf(
		"BigEndian: 0x%x\n", binary.BigEndian.Uint32(bytes),
	)
}
