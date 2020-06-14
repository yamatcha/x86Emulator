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

func movRm32Imm32(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	value := getCode32(emu, 0)
	emu.Eip += 4
	setRm32(emu, &modrm, value)
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

func movRm8R8(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	r8 := getR8(emu, &modrm)
	setRm8(emu, &modrm, r8)
}

func movR8Imm8(emu *Emulator) {
	reg := GetCode8(emu, 0) - 0xb0
	emu.Registers.setRegister8(reg, GetCode8(emu, 1))
	emu.Eip += 2
}

func movR8Rm8(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	rm8 := getRm8(emu, &modrm)
	setR8(emu, &modrm, byte(rm8))
}

func addRm32Imm8(emu *Emulator, modrm *ModRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := int32(getSignCode8(emu, 0))
	emu.Eip++
	setRm32(emu, modrm, rm32+uint32(imm8))
}

func addRm32R32(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	r32 := getR32(emu, &modrm)
	rm32 := getRm32(emu, &modrm)
	setRm32(emu, &modrm, rm32+r32)
}

func subRm32Imm8(emu *Emulator, modrm *ModRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := uint32(getSignCode8(emu, 0))
	emu.Eip++
	result := uint64(rm32) - uint64(imm8)
	setRm32(emu, modrm, rm32-imm8)
	updateEflagsSub(emu, rm32, imm8, result)
}

func cmpAlImm8(emu *Emulator) {
	value := uint32(GetCode8(emu, 1))
	al := uint32(emu.Registers.getRegister8(AL))
	result := uint64(al) - uint64(value)
	updateEflagsSub(emu, al, value, result)
	emu.Eip += 2
}

func cmpEaxImm32(emu *Emulator) {
	value := getCode32(emu, 1)
	eax := emu.Registers.GetRegister32(EAX)
	result := uint64(eax) - uint64(value)
	updateEflagsSub(emu, eax, value, result)
	emu.Eip++
}

func cmpR32Rm32(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	r32 := getR32(emu, &modrm)
	rm32 := getRm32(emu, &modrm)
	result := uint64(r32) - uint64(rm32)
	updateEflagsSub(emu, r32, rm32, result)
}

func cmpRm32Imm8(emu *Emulator, modrm *ModRM) {
	rm32 := getRm32(emu, modrm)
	imm8 := int32(getSignCode8(emu, 0))
	result := uint64(rm32) - uint64(imm8)
	updateEflagsSub(emu, rm32, uint32(imm8), result)
}

func incRm32(emu *Emulator, modrm *ModRM) {
	value := getRm32(emu, modrm)
	setRm32(emu, modrm, value+1)
}

func incR32(emu *Emulator) {
	reg := GetCode8(emu, 0) - 0x40
	emu.Registers.setRegister32(reg, emu.Registers.GetRegister32(reg)+1)
	emu.Eip++
}

func code83(emu *Emulator) {
	emu.Eip++
	var modrm ModRM
	parseModrm(emu, &modrm)
	switch modrm.opecode {
	case 0:
		addRm32Imm8(emu, &modrm)
	case 5:
		subRm32Imm8(emu, &modrm)
	case 7:
		cmpRm32Imm8(emu, &modrm)
	default:
		fmt.Println("not implemented: 83", modrm.opecode)
		os.Exit(1)
	}
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

func jc(emu *Emulator) {
	if isCarry(emu) {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2
	} else {
		emu.Eip += 2
	}
}

func jnc(emu *Emulator) {
	if isCarry(emu) {
		emu.Eip += 2
	} else {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2

	}
}

func jz(emu *Emulator) {
	if isZero(emu) {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2
	} else {
		emu.Eip += 2
	}
}

func jnz(emu *Emulator) {
	if isZero(emu) {
		emu.Eip += 2
	} else {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2

	}
}

func js(emu *Emulator) {
	if isSign(emu) {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2
	} else {
		emu.Eip += 2
	}
}

func jns(emu *Emulator) {
	if isSign(emu) {
		emu.Eip += 2
	} else {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2

	}
}

func jo(emu *Emulator) {
	if isOverflow(emu) {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2
	} else {
		emu.Eip += 2
	}
}

func jno(emu *Emulator) {
	if isOverflow(emu) {
		emu.Eip += 2
	} else {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2

	}
}

func jl(emu *Emulator) {
	if isSign(emu) != isOverflow(emu) {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2
	} else {
		emu.Eip += 2
	}
}

func jle(emu *Emulator) {
	if isZero(emu) || isSign(emu) != isOverflow(emu) {
		emu.Eip += uint32(getSignCode8(emu, 1)) + 2
	} else {
		emu.Eip += 2
	}
}

func pushImm8(emu *Emulator) {
	value := GetCode8(emu, 1)
	push32(emu, uint32(value))
	emu.Eip += 2
}

func pushImm32(emu *Emulator) {
	value := getCode32(emu, 1)
	push32(emu, value)
	emu.Eip += 5
}

func pushR32(emu *Emulator) {
	reg := GetCode8(emu, 0) - 0x50
	push32(emu, emu.Registers.GetRegister32(reg))
	emu.Eip++
}

func popR32(emu *Emulator) {
	reg := GetCode8(emu, 0) - 0x58
	emu.Registers.setRegister32(reg, pop32(emu))
	emu.Eip++
}

func callRel32(emu *Emulator) {
	diff := getSignCode32(emu, 1)
	push32(emu, emu.Eip+5)
	emu.Eip += uint32(diff + 5)
}

func leave(emu *Emulator) {
	ebp := emu.Registers.GetRegister32(EBP)
	emu.Registers.setRegister32(ESP, ebp)
	emu.Registers.setRegister32(EBP, pop32(emu))
	emu.Eip++
}

func ret(emu *Emulator) {
	emu.Eip = pop32(emu)
}

func nearJump(emu *Emulator) {
	diff := getSignCode32(emu, 1)
	emu.Eip += uint32(diff + 5)
}

func shortJump(emu *Emulator) {
	var diff int8 = getSignCode8(emu, 1)
	emu.Eip += uint32(diff + 2)
}

func inAlDx(emu *Emulator) {
	address := emu.Registers.GetRegister32(EDX) & 0xffff
	value := ioIn8(uint16(address))
	emu.Registers.setRegister8(AL, value)
	emu.Eip++
}

func outDxAl(emu *Emulator) {
	address := emu.Registers.GetRegister32(EDX) & 0xffff
	value := emu.Registers.getRegister8(AL)
	ioOut8(uint16(address), value)
	emu.Eip++
}

func swi(emu *Emulator) {
	intIndex := GetCode8(emu, 1)
	emu.Eip += 2
	switch intIndex {
	case 0x10:
		biosVideo(emu)
	default:
		fmt.Printf("unknown interrupt: 0x%02x\n", intIndex)
	}
}

func Instructions(index byte) (func(emu *Emulator), error) {
	switch {
	case 0x01 == index:
		return addRm32R32, nil
	case 0x3b == index:
		return cmpR32Rm32, nil
	case 0x3c == index:
		return cmpAlImm8, nil
	case 0x3d == index:
		return cmpEaxImm32, nil
	case 0x40 <= index && index < 0x40+8:
		return incR32, nil
	case 0x50 <= index && index < 0x50+8:
		return pushR32, nil
	case 0x58 <= index && index < 0x58+8:
		return popR32, nil
	case 0x68 == index:
		return pushImm32, nil
	case 0x6a == index:
		return pushImm8, nil

	case 0x70 == index:
		return jo, nil
	case 0x71 == index:
		return jno, nil
	case 0x72 == index:
		return jc, nil
	case 0x73 == index:
		return jnc, nil
	case 0x74 == index:
		return jz, nil
	case 0x75 == index:
		return jnz, nil
	case 0x78 == index:
		return js, nil
	case 0x79 == index:
		return jns, nil
	case 0x7c == index:
		return jl, nil
	case 0x7e == index:
		return jle, nil

	case 0x83 == index:
		return code83, nil
	case 0x88 == index:
		return movRm8R8, nil
	case 0x89 == index:
		return movRm32R32, nil
	case 0x8a == index:
		return movR8Rm8, nil
	case 0x8b == index:
		return movR32Rm32, nil
	case 0xb0 <= index && index < 0xb0+8:
		return movR8Imm8, nil
	case 0xb8 <= index && index < 0xb8+8:
		return movR32Imm32, nil
	case 0xc3 == index:
		return ret, nil
	case 0xc7 == index:
		return movRm32Imm32, nil
	case 0xc9 == index:
		return leave, nil

	case 0xcd == index:
		return swi, nil

	case 0xe8 == index:
		return callRel32, nil
	case 0xe9 == index:
		return nearJump, nil
	case 0xeb == index:
		return shortJump, nil
	case 0xec == index:
		return inAlDx, nil
	case 0xee == index:
		return outDxAl, nil
	case 0xff == index:
		return codeff, nil
	}
	return nil, fmt.Errorf("Error: Invalid index of instructions %x", index)
}
