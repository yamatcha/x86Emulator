package emulator

import "fmt"

var biosToTerminal []int = []int{
	30, 34, 32, 36, 31, 35, 33, 37}

func putString(s string) {
	for _, i := range s {
		ioOut8(0x03f8, uint8(i))
	}
}

func biosVideoTeletype(emu *Emulator) {
	color := emu.Registers.getRegister8(BL) & 0x0f
	ch := emu.Registers.getRegister8(AL)
	var buf string
	var bright int
	terminalColor := biosToTerminal[color&0x07]
	if color&0x08 == 1 {
		bright = 1
	} else {
		bright = 0
	}
	buf = fmt.Sprintf("\x1b[%d;%dm%c\x1b[0m", bright, terminalColor, ch)
	putString(buf)
}

func biosVideo(emu *Emulator) {
	function := emu.Registers.getRegister8(AH)
	switch function {
	case 0x0e:
		biosVideoTeletype(emu)
	default:
		fmt.Printf("not implemented BIOS video function 0x%02x\n", function)
	}
}
