package isa

import "fmt"

type Instruction struct {
	Opcode Opcode
	Rd     uint8
	Rs     uint8
	Rx     uint8
	Func2  uint8
	Func3  uint8
	Func5  uint8
	Imm    uint16
	Ex     uint8
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

func (i Instruction) String() string {
	switch i.Opcode.Opc {
	case OP_ALU_REG_REG:
		if name, exists := mapFunc5ToAluRegReg[i.Func5]; exists {
			rdBanked := i.Rd
			if (i.Ex & 0x01) != 0 {
				rdBanked += 8
			}
			rsBanked := i.Rs
			if (i.Ex & 0x02) != 0 {
				rsBanked += 8
			}
			rsStr := fmt.Sprintf("R%d", rdBanked)
			rdStr := fmt.Sprintf("R%d", rsBanked)
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
			return fmt.Sprintf("%s %s, %s", name, rdStr, rsStr)
		}
	case OP_LD_REG_IMM:
		return fmt.Sprintf("LD.%d R%d, 0x%x", i.Func2, i.Rd, i.Imm)
	}
	return ""
}
