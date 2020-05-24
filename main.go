package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yamatcha/x86Emulator/emulator"
)

func dumpRegisters(emu *emulator.Emulator) {
	for i := 0; i < emulator.RegistersCount; i++ {
		fmt.Printf("%s = %08x\n", emulator.RegistersName[i], emu.Registers.GetRegister32(byte(i)))
	}

	fmt.Printf("EIP = %08x\n", emu.Eip)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("few arguments")
		os.Exit(1)
	}
	filename := os.Args[1]
	emu := emulator.NewEmulator(emulator.MemorySize, 0x7c00, 0x7c00)
	if emu == nil {
		fmt.Println("Error: Value of eip or esp is invalid")
		os.Exit(1)
	}
	binary, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer binary.Close()
	mem, err := ioutil.ReadAll(binary)
	emu.Memory = append(emu.Memory[:0x7c00], mem...)
	for emu.Eip < emulator.MemorySize {
		var code byte = emulator.GetCode8(emu, 0)
		fmt.Printf("EIP = %X, Code = %02X\n", emu.Eip, code)
		operator, err := emulator.Instructions(code)
		if err != nil {
			panic(err)
		}
		operator(emu)
		if emu.Eip == 0x00 {
			fmt.Printf("\n\n End of program.\n\n")
			break
		}
	}
	dumpRegisters(emu)
}
