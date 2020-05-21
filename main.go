package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type registerID int

const (
	registersCount = 8
	memorySize     = 1024 * 1024
)

type registers struct {
	eax uint32
	ecx uint32
	edx uint32
	ebx uint32
	esp uint32
	ebp uint32
	esi uint32
	edi uint32
}

var registersName []string = []string{
	"EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI"}

type emulator struct {
	registers *registers
	eflags    uint32
	memory    []byte
	eip       uint32
}

func newEmulator(size uint32, eip uint32, esp uint32) *emulator {
	if eip < 0 || esp < 0 {
		return nil
	}
	r := registers{esp: esp}
	e := emulator{eip: eip, memory: make([]byte, size)}
	e.registers = &r
	return &e
}

func dumpRegisters(emu *emulator) {
	for i := 0; i < registersCount; i++ {
		fmt.Printf("%s = %08x\n", registersName[i], emu.registers.getRegisterValue(byte(i)))
	}

	fmt.Printf("EIP = %08x\n", emu.eip)
}

func (reg *registers) setRegisterValue(index byte, value uint32) {
	switch index {
	case 0:
		reg.eax = value
	case 1:
		reg.ecx = value
	case 2:
		reg.edx = value
	case 3:
		reg.ebx = value
	case 4:
		reg.esp = value
	case 5:
		reg.ebp = value
	case 6:
		reg.esi = value
	case 7:
		reg.edi = value
	}
	return
}

func (reg *registers) getRegisterValue(index byte) uint32 {
	switch index {
	case 0:
		return reg.eax
	case 1:
		return reg.ecx
	case 2:
		return reg.edx
	case 3:
		return reg.ebx
	case 4:
		return reg.esp
	case 5:
		return reg.ebp
	case 6:
		return reg.esi
	case 7:
		return reg.edi
	}
	return 0
}

func getCode8(emu *emulator, index int) byte {
	return emu.memory[int(emu.eip)+index]
}

func getSignCode8(emu *emulator, index int) int8 {
	return int8(emu.memory[int(emu.eip)+index])
}

func getCode32(emu *emulator, index int) uint32 {
	var ret uint32 = 0
	for i := 0; i < 4; i++ {
		ret |= uint32(getCode8(emu, index+i) << (i * 8))
	}
	return ret
}

func movR32Imm32(emu *emulator) {
	var reg byte = getCode8(emu, 0) - 0xB8
	var value uint32 = getCode32(emu, 1)
	emu.registers.setRegisterValue(reg, value)
	emu.eip += 5
}

func shortJump(emu *emulator) {
	var diff int8 = getSignCode8(emu, 1)
	emu.eip += uint32(diff + 2)
}

func instructions(index byte) (func(emu *emulator), error) {
	switch {
	case 0xb7 < index && index < 0xb8+8:
		return movR32Imm32, nil
	case 0xeb == index:
		return shortJump, nil
	}
	return nil, errors.New("Error: Invalid index of instructions")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("few arguments")
		os.Exit(1)
	}
	filename := os.Args[1]
	emu := newEmulator(memorySize, 0x0000, 0x7c00)
	if emu == nil {
		fmt.Println("Error: Value of eip or esp is invalid")
		os.Exit(1)
	}
	binary, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer binary.Close()
	emu.memory, err = ioutil.ReadAll(binary)
	for emu.eip < memorySize {
		var code byte = getCode8(emu, 0)
		fmt.Printf("EIP = %X, Code = %02X\n", emu.eip, code)
		operator, err := instructions(code)
		if err != nil {
			panic(err)
		}
		operator(emu)
		if emu.eip == 0x00 {
			fmt.Printf("\n\n End of program.\n\n")
			break
		}
	}
	dumpRegisters(emu)
}
