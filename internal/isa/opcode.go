package isa

type Opcode struct {
	Opc  uint8
	Name string
	Type uint8
}

func (o Opcode) String() string {
	return o.Name
}

const (
	OP_ALU_REG_REG       = 0x00
	OP_LD_REG_IMM        = 0x01
	OP_ALU_REG_IMM       = 0x02
	OP_REG_MEM           = 0x03
	OP_REG_MEM_IMM       = 0x04
	OP_JUMP_CALL_RET_REG = 0x05
	OP_JUMP_REL          = 0x06
	OP_CALL_REL          = 0x07
)

const (
	OP_TYPE_2_REG     = 0x00
	OP_TYPE_3_REG     = 0x01
	OP_TYPE_7_IMM_REG = 0x02
	OP_TYPE_8_IMM_REG = 0x03
	OP_TYPE_13_IMM    = 0x04
)

var (
	OpAluRegReg      = Opcode{Opc: OP_ALU_REG_REG, Name: "ALU_REG_REG", Type: OP_TYPE_2_REG}
	OpLdRegImm       = Opcode{Opc: OP_LD_REG_IMM, Name: "LD_REG_IMM", Type: OP_TYPE_8_IMM_REG}
	OpAluRegImm      = Opcode{Opc: OP_ALU_REG_IMM, Name: "ALU_REG_IMM", Type: OP_TYPE_7_IMM_REG}
	OpRegMem         = Opcode{Opc: OP_REG_MEM, Name: "REG_MEM", Type: OP_TYPE_2_REG}
	OpRegMemImm      = Opcode{Opc: OP_REG_MEM_IMM, Name: "REG_MEM_IMM", Type: OP_TYPE_3_REG}
	OpJumpCallRetReg = Opcode{Opc: OP_JUMP_CALL_RET_REG, Name: "JUMP_CALL_RET_REG", Type: OP_TYPE_3_REG}
	OpJumpRel        = Opcode{Opc: OP_JUMP_REL, Name: "JUMP_REL", Type: OP_TYPE_3_REG}
	OpCallRel        = Opcode{Opc: OP_CALL_REL, Name: "CALL_REL", Type: OP_TYPE_13_IMM}
)

var opcodeMap = map[uint8]Opcode{
	OP_ALU_REG_REG:       OpAluRegReg,
	OP_LD_REG_IMM:        OpLdRegImm,
	OP_ALU_REG_IMM:       OpAluRegImm,
	OP_REG_MEM:           OpRegMem,
	OP_REG_MEM_IMM:       OpRegMemImm,
	OP_JUMP_CALL_RET_REG: OpJumpCallRetReg,
	OP_JUMP_REL:          OpJumpRel,
	OP_CALL_REL:          OpCallRel,
}

func GetOpcode(opc uint8) Opcode {
	opcode, _ := opcodeMap[opc]
	return opcode
}
