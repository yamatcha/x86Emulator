package emulator

type RegisterID int

const (
	RegistersCount = 8
	MemorySize     = 1024 * 1024
	ESP            = 4
	EBP            = 5
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

func setMomory8(emu *Emulator, address uint32, value uint32) {
	emu.Memory[address] = byte(value & 0xff)
}

func setMemory32(emu *Emulator, address uint32, value uint32) {
	for i := uint32(0); i < 4; i++ {
		setMomory8(emu, address+i, value>>(i*8))
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
