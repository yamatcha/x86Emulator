package emulator

import (
	"fmt"
	"os"
)

func movR32Imm32(emu *Emulator) {
	reg := GetCode8(emu, 0) - 0xB8
	value := getCode32(emu, 1)
	emu.Registers.setRegister32(reg, value)
	emu.Eip += 5
}

func movRm32R32(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	r32 := getR32(emu, &modrm)
	setRm32(emu, &modrm, r32)
}

func movR32Rm32(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	rm32 := getRm32(emu, &modrm)
	setR32(emu, &modrm, rm32)
}

func addRm32R32(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	r32 := getR32(emu, &modrm)
	rm32 := getRm32(emu, &modrm)
	setRm32(emu, &modrm, rm32+r32)
}

func movRm32Imm32(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	value := getCode32(emu, 0)
	emu.Eip += 4
	setRm32(emu, &modrm, value)
}

func subRm32Imm8(emu *Emulator, modrm *ModRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := uint32(getSignCode8(emu, 0))
	emu.Eip++
	setRm32(emu, modrm, rm32-imm8)
}

func code83(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	switch modrm.opecode {
	case 5:
		subRm32Imm8(emu, &modrm)
	default:
		fmt.Println("not implemented: 83", modrm.opecode)
		os.Exit(1)
	}
}

func incRm32(emu *Emulator, modrm *ModRM) {
	value := getRm32(emu, modrm)
	setRm32(emu, modrm, value+1)
}

func codeff(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	switch modrm.opecode {
	case 0:
		incRm32(emu, &modrm)
	default:
		fmt.Println("not implemented: FF ", modrm.opecode)
		os.Exit(1)
	}
}

func nearJump(emu *Emulator) {
	diff := getSignCode32(emu, 1)
	emu.Eip += uint32(diff + 5)
}

func shortJump(emu *Emulator) {
	var diff int8 = getSignCode8(emu, 1)
	emu.Eip += uint32(diff + 2)
}

func Instructions(index byte) (func(emu *Emulator), error) {
	switch {
	case 0x01 == index:
		return addRm32R32, nil
	case 0x83 == index:
		return code83, nil
	case 0x89 == index:
		return movRm32R32, nil
	case 0x8b == index:
		return movR32Rm32, nil
	case 0xb7 < index && index < 0xb8+8:
		return movR32Imm32, nil
	case 0xc7 == index:
		return movRm32Imm32, nil
	case 0xe9 == index:
		return nearJump, nil
	case 0xeb == index:
		return shortJump, nil
	case 0xff == index:
		return codeff, nil
	}
	return nil, fmt.Errorf("Error: Invalid index of instructions %d", index)
}
