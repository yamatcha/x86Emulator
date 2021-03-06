package emulator

type RegisterID int

const (
	RegistersCount        = 8
	MemorySize            = 1024 * 1024
	AL                    = 0
	BL                    = 3
	AH                    = 4
	EAX                   = 0
	EDX                   = 2
	ESP                   = 4
	EBP                   = 5
	CARRYFLAG      uint32 = 1
	ZEROFLAG       uint32 = 1 << 6
	SIGNFLAG       uint32 = 1 << 7
	OVERFLOWFLAG   uint32 = 1 << 11
)

type Registers struct {
	eax uint32
	ecx uint32
	edx uint32
	ebx uint32
	esp uint32
	ebp uint32
	esi uint32
	edi uint32
}

var RegistersName []string = []string{
	"EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI"}

type Emulator struct {
	Registers *Registers
	eflags    uint32
	Memory    []byte
	Eip       uint32
}

func NewEmulator(size uint32, eip uint32, esp uint32) *Emulator {
	if eip < 0 || esp < 0 {
		return nil
	}
	r := Registers{esp: esp}
	e := Emulator{Eip: eip, Memory: make([]byte, size)}
	e.Registers = &r
	return &e
}

func (reg *Registers) setRegister32(index byte, value uint32) {
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

func (reg *Registers) GetRegister32(index byte) uint32 {
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

func (reg *Registers) getRegister8(index byte) byte {
	if index < 4 {
		return byte(reg.GetRegister32(index) & 0xff)
	} else {
		return byte((reg.GetRegister32(index-4) >> 8) & 0xff)
	}
}

func (reg *Registers) setRegister8(index byte, value byte) {
	if index < 4 {
		r := reg.GetRegister32(index) & 0xffffff00
		reg.setRegister32(index, r|uint32(value))
	} else {
		r := reg.GetRegister32(index-4) & 0xffff00ff
		reg.setRegister32(index-4, r|uint32(value)<<8)
	}
}

func GetCode8(emu *Emulator, index int) byte {
	return emu.Memory[int(emu.Eip)+index]
}

func getSignCode8(emu *Emulator, index int) int8 {
	return int8(emu.Memory[int(emu.Eip)+index])
}

func getCode32(emu *Emulator, index int) uint32 {
	var ret uint32 = 0
	for i := 0; i < 4; i++ {
		ret |= uint32(GetCode8(emu, index+i)) << (i * 8)
	}
	return ret
}

func getSignCode32(emu *Emulator, index int) int32 {
	return int32(getCode32(emu, index))
}

func setMemory8(emu *Emulator, address uint32, value uint32) {
	emu.Memory[address] = byte(value & 0xff)
}

func setMemory32(emu *Emulator, address uint32, value uint32) {
	for i := uint32(0); i < 4; i++ {
		setMemory8(emu, address+i, value>>(i*8))
	}
}

func getMemory8(emu *Emulator, address uint32) uint32 {
	return uint32(emu.Memory[address])
}

func getMemory32(emu *Emulator, address uint32) uint32 {
	ret := uint32(0)
	for i := uint32(0); i < 4; i++ {
		ret |= getMemory8(emu, address+i) << (8 * i)
	}
	return ret
}

func push32(emu *Emulator, value uint32) {
	address := emu.Registers.GetRegister32(ESP) - 4
	emu.Registers.setRegister32(ESP, address)
	setMemory32(emu, address, value)
}

func pop32(emu *Emulator) uint32 {
	address := emu.Registers.GetRegister32(ESP)
	ret := getMemory32(emu, address)
	emu.Registers.setRegister32(ESP, address+4)
	return ret
}

func updateEflagsSub(emu *Emulator, v1, v2 uint32, result uint64) {
	sign1 := int(v1 >> 31)
	sign2 := int(v2 >> 31)
	signr := int((result >> 31) & 1)
	setCarry(emu, result>>32 != 0)
	setZero(emu, result == 0)
	setSign(emu, signr != 0)
	setOverflow(emu, sign1 != sign2 && sign1 != signr)
}

func setCarry(emu *Emulator, isCarry bool) {
	if isCarry {
		emu.eflags |= CARRYFLAG
	} else {
		emu.eflags &= ^CARRYFLAG
	}
}

func setZero(emu *Emulator, isZero bool) {
	if isZero {
		emu.eflags |= ZEROFLAG
	} else {
		emu.eflags &= ^ZEROFLAG
	}
}

func setSign(emu *Emulator, isSign bool) {
	if isSign {
		emu.eflags |= SIGNFLAG
	} else {
		emu.eflags &= ^SIGNFLAG
	}
}

func setOverflow(emu *Emulator, isOverflow bool) {
	if isOverflow {
		emu.eflags |= OVERFLOWFLAG
	} else {
		emu.eflags &= ^OVERFLOWFLAG
	}
}

func isCarry(emu *Emulator) bool {
	return (emu.eflags & CARRYFLAG) != 0
}

func isZero(emu *Emulator) bool {
	return (emu.eflags & ZEROFLAG) != 0
}
func isSign(emu *Emulator) bool {
	return (emu.eflags & SIGNFLAG) != 0
}
func isOverflow(emu *Emulator) bool {
	return (emu.eflags & OVERFLOWFLAG) != 0
}
