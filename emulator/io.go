package emulator

import (
	"bufio"
	"fmt"
	"os"
)

func ioIn8(address uint16) uint8 {
	switch address {
	case 0x03f8:
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadByte()
		return input
	default:
		return 0
	}
}

func ioOut8(address uint16, value uint8) {
	switch address {
	case 0x03f8:
		fmt.Printf("%c", value)
	}
}
