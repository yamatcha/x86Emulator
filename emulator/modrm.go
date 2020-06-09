package emulator

import (
	"fmt"
	"os"
)

type ModRM struct {
	mod     byte
	opecode byte
	// regIndex byte
	rm     byte
	sib    byte
	disp32 uint32
	disp8  int8
}

func parseModrm(emu *Emulator, modrm *ModRM) {
	var code byte

	code = GetCode8(emu, 0)
	modrm.mod = ((code & 0xc0) >> 6)
	modrm.opecode = ((code & 0x38) >> 3)
	// modrm.regIndex = modrm.opecode
	modrm.rm = code & 0x07

	emu.Eip++
	if modrm.mod != 3 && modrm.rm == 4 {
		modrm.sib = GetCode8(emu, 0)
		emu.Eip++
	}

	if (modrm.mod == 0 && modrm.rm == 5) || modrm.mod == 2 {
		modrm.disp32 = getCode32(emu, 0)
		emu.Eip += 4
	} else if modrm.mod == 1 {
		modrm.disp8 = getSignCode8(emu, 0)
		emu.Eip++
	}
}

func calcMemoryAddress(emu *Emulator, modrm *ModRM) uint32 {
	if modrm.mod == 0 {
		if modrm.rm == 4 {
			fmt.Println("not implemented Mod mod = 0, rm = 4")
			os.Exit(0)
		} else if modrm.mod == 5 {
			return modrm.disp32
		} else {
			return emu.Registers.GetRegister32(modrm.rm)
		}
	} else if modrm.mod == 1 {
		if modrm.rm == 4 {
			fmt.Println("not implemented Mod mod = 1, rm = 4")
			os.Exit(0)
		} else {
			return emu.Registers.GetRegister32(modrm.rm) + uint32(modrm.disp8)
		}
	} else if modrm.mod == 2 {
		if modrm.rm == 4 {
			fmt.Println("not implemented Mod mod = 2, rm = 4")
			os.Exit(0)
		} else {
			return emu.Registers.GetRegister32(modrm.rm) + modrm.disp32
		}
	} else {
		fmt.Println("not implemented Mod mod = 3")
		os.Exit(0)
	}
	return 0
}

func setRm32(emu *Emulator, modrm *ModRM, value uint32) {
	if modrm.mod == 3 {
		emu.Registers.setRegister32(modrm.rm, value)
	} else {
		address := calcMemoryAddress(emu, modrm)
		setMemory32(emu, address, value)
	}
}

func getRm32(emu *Emulator, modrm *ModRM) uint32 {
	if modrm.mod == 3 {
		return emu.Registers.GetRegister32(modrm.rm)
	} else {
		address := calcMemoryAddress(emu, modrm)
		return getMemory32(emu, address)
	}
}

func setR32(emu *Emulator, modrm *ModRM, value uint32) {
	emu.Registers.setRegister32(modrm.opecode, value)
	// emu.Registers.setRegister32(modrm.regIndex, value)
}

func getR32(emu *Emulator, modrm *ModRM) uint32 {
	return emu.Registers.GetRegister32(modrm.opecode)
	// return emu.Registers.GetRegister32(modrm.regIndex)
}

func getRm8(emu *Emulator, modrm *ModRM) byte {
	if modrm.mod == 3 {
		return emu.Registers.getRegister8(modrm.rm)
	} else {
		address := calcMemoryAddress(emu, modrm)
		return byte(getMemory8(emu, address))
	}
}

func setRm8(emu *Emulator, modrm *ModRM, value byte) {
	if modrm.mod == 3 {
		emu.Registers.setRegister8(modrm.rm, value)
	} else {
		address := calcMemoryAddress(emu, modrm)
		setMemory8(emu, address, uint32(value))
	}
}

func getR8(emu *Emulator, modrm *ModRM) byte {
	return emu.Registers.getRegister8(modrm.opecode)
}

func setR8(emu *Emulator, modrm *ModRM, value byte) {
	emu.Registers.setRegister8(modrm.opecode, value)
}
