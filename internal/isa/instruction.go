package isa

import "fmt"

type Instruction struct {
	Opcode  Opcode
	Rd      uint8
	Rs      uint8
	Rx      uint8
	Func2   uint8
	Func3   uint8
	Func5   uint8
	Imm     int16
	Ex      uint8
	Address uint32
	Label   string
}

func (i Instruction) Pack() uint16 {
	var code uint16 = uint16(i.Opcode.Opc)

	//Pack instruction by opcode type
	switch i.Opcode.Type {
	case OP_TYPE_2_REG:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Rs) << 13
		code |= uint16(i.Ex) << 11
		code |= uint16(i.Func5) << 6
	case OP_TYPE_8_IMM_REG:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Func2) << 6
		code |= uint16(i.Imm) << 8
	case OP_TYPE_7_IMM_REG:
		code |= uint16(i.Rd) << 3
		code |= uint16(i.Func3) << 6
		code |= uint16(i.Imm) << 9
		//TODO:
	}

	return code
}

func (i Instruction) makeEx(bankd, banks uint8) Instruction {
	i.Ex = (banks << 1) | bankd
	return i
}

func getRd(inst uint16) uint8 {
	return uint8(inst>>3) & 0b0000_0111
}

func getRs(inst uint16) uint8 {
	return uint8(inst>>13) & 0b0000_0111
}

func getFunc5(inst uint16) uint8 {
	return uint8(inst>>6) & 0b0001_1111
}

func getFunc3(inst uint16) uint8 {
	return uint8(inst>>6) & 0b0000_0111
}

func getFunc2(inst uint16) uint8 {
	return uint8(inst>>6) & 0b0000_0011
}

func getEx(inst uint16) uint8 {
	return uint8(inst>>11) & 0b0000_0011
}

func getImm8(inst uint16) int16 {
	return int16(int8(inst >> 8))
}

func getImm7(inst uint16) int16 {
	return int16(int8(inst>>8) >> 1)
}

func Parse(inst uint16) Instruction {
	i := Instruction{}
	i.Opcode = GetOpcode(uint8(inst) & 0b0000_0111)
	switch i.Opcode.Type {
	case OP_TYPE_3_REG:
		i.Rd = getRd(inst)
		i.Rs = getRs(inst)
		i.Func5 = getFunc5(inst)
		i.Ex = getEx(inst)
	case OP_TYPE_8_IMM_REG:
		i.Rd = getRd(inst)
		i.Func2 = getFunc2(inst)
		i.Imm = int16(getImm8(inst))
	case OP_TYPE_7_IMM_REG:
		i.Rd = getRd(inst)
		i.Func3 = getFunc3(inst)
		i.Imm = int16(getImm7(inst))
		//TODO:
	}
	return i
}

var mapFunc5ToAluRegReg = map[uint8]string{
	0:  "LD",
	1:  "ADD",
	2:  "SUB",
	3:  "SHL",
	4:  "SHR",
	5:  "AND",
	6:  "OR",
	7:  "XOR",
	8:  "NOT",
	9:  "MUL",
	10: "LD",
	11: "LD",
	12: "LD",
	13: "LD",
	14: "INT",
}

var mapAluToFunc5 = map[string]uint8{
	"ADD": 1,
	"SUB": 2,
	"SHL": 3,
	"SHR": 4,
	"AND": 5,
	"OR":  6,
	"XOR": 7,
	"NOT": 8,
	"MUL": 9,
}

func getFunc5FromALU(alu string) (uint8, error) {
	func5, exists := mapAluToFunc5[alu]
	if !exists {
		return 0, fmt.Errorf("invalid ALU name")
	}
	return func5, nil
}

func (i Instruction) String() string {
	switch i.Opcode.Opc {
	case OP_ALU_REG_REG:
		if name, exists := mapFunc5ToAluRegReg[i.Func5]; exists {
			rdBanked := i.Rd
			if (i.Ex & 0x01) != 0 {
				rdBanked += 7
			}
			rsBanked := i.Rs
			if (i.Ex & 0x02) != 0 {
				rsBanked += 7
			}
			rdStr := fmt.Sprintf("R%d", rdBanked)
			rsStr := fmt.Sprintf("R%d", rsBanked)
			switch i.Func5 {
			case 10:
				rsStr = "SP"
			case 11:
				rsStr = "LR"
			case 12:
				rdStr = "SP"
			case 13:
				rdStr = "LR"
			}
			return fmt.Sprintf("%s\t%s, %s", name, rdStr, rsStr)
		}
	case OP_LD_REG_IMM:
		rdBanked := i.Rd
		if (i.Func2 & 0b0000_0010) != 0 {
			rdBanked += 7
		}
		return fmt.Sprintf("LD.%d\tR%d, 0x%02X(%d)", i.Func2&0b0000_0001, rdBanked, uint8(i.Imm), i.Imm)
	case OP_ALU_REG_IMM:
		rdBanked := i.Rd
		if (i.Func3 & 0b0000_0100) != 0 {
			rdBanked += 7
		}
		var operand string
		funcOper := i.Func3 & 0b0000_0011
		imm := i.Imm
		switch funcOper {
		case 0:
			operand = "ADD"
			if i.Imm < 0 {
				operand = "SUB"
				imm = ^imm + 1
			}
		case 1:
			operand = "SHL"
			if i.Imm < 0 {
				operand = "SHR"
				imm = ^imm + 1
			}
		case 2:
			operand = "LDI"
		case 3:
			operand = "DJNZ"
		}
		switch funcOper {
		case 2: //LDI
			return fmt.Sprintf("%s\tR%d, 0x%02X(%d) -> %08X", operand, rdBanked, uint8(i.Imm), i.Imm, uint32(int32(i.Address)+int32((-i.Imm)<<1)))
		case 3: //DJNZ
			return fmt.Sprintf("%s\tR%d, 0x%02X(%d) -> %08X", operand, rdBanked, uint8(i.Imm), i.Imm, uint32(int32(i.Address)+int32(i.Imm<<1)))
		default:
			return fmt.Sprintf("%s\tR%d, 0x%02X(%d)", operand, rdBanked, uint8(imm), imm)
		}
		//TODO:
	}
	return ""
}
