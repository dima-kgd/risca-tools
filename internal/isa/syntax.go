package isa

import "fmt"

const (
	RuleALURegReg = iota
	LdRegImm
	AluRegImm
	Label
)

type Rule struct {
	Type      uint8
	Syntax    [][]uint8
	Opcode    Opcode
	ParseFunc func(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error)
}

var aluRegRegSyntax = [][]uint8{{TK_ALU, TK_LD}, {TK_REG, TK_REG_SP, TK_REG_LR}, {TK_COMMA}, {TK_REG, TK_REG_SP, TK_REG_LR}}
var ldRegImmSyntax = [][]uint8{{TK_LD_0, TK_LD_1}, {TK_REG}, {TK_COMMA}, {TK_NUMBER}}
var labelSyntax = [][]uint8{{TK_LABEL}, {TK_COLON}}
var aluRegImmSyntax = [][]uint8{{TK_ALU, TK_LDI, TK_DJNZ}, {TK_REG}, {TK_COMMA}, {TK_NUMBER, TK_LABEL}}

var syntaxRules = []Rule{
	{Type: RuleALURegReg, Syntax: aluRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG), ParseFunc: ParseAluRegReg},
	{Type: LdRegImm, Syntax: ldRegImmSyntax, Opcode: GetOpcode(OP_LD_REG_IMM), ParseFunc: ParseLdRegImm},
	{Type: AluRegImm, Syntax: aluRegImmSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: ParseAluRegImm},
	{Type: Label, Syntax: labelSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: ParseLabel},
}

func parseRegister(tokenrd string) (uint8, uint8, error) {
	bankd := uint8(0)
	rd, found := GetRegisterNumber(tokenrd)
	if !found {
		return 0, 0, fmt.Errorf("Invalid register: %s", tokenrd)
	}
	if rd > 7 {
		rd = rd - 7
		bankd = 1
	}
	return rd, bankd, nil
}

func ParseAluRegReg(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, bankd, rs, banks uint8
	var err error
	var func5 uint8

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	regST := tokens[tokenpos+3]
	instr := Instruction{}

	if regDT.T != TK_REG_LR && regDT.T != TK_REG_SP {
		rd, bankd, err = parseRegister(regDT.Tk)
		if err != nil {
			return tokens, err
		}
	}

	if regST.T != TK_REG_LR && regST.T != TK_REG_SP {
		rs, banks, err = parseRegister(regST.Tk)
		if err != nil {
			return tokens, err
		}
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Rs = rs
	instr = instr.makeEx(bankd, banks)
	instr.Address = parser.CurAddress
	func5 = 0 //LD for default LD REG, REG

	if aluT.T == TK_LD { // LD instruction
		if regDT.T == TK_REG_LR { // LD LR, REG
			if regST.T != TK_REG {
				return tokens, fmt.Errorf("Source should be R0-R15, found %v", regST)
			}
			func5 = 13
		} else if regDT.T == TK_REG_SP { //LD SP, REG
			if regST.T != TK_REG {
				return tokens, fmt.Errorf("Source should be R0-R15, found %v", regST)
			}
			func5 = 12
		} else if regST.T == TK_REG_LR { //LD REG, LR
			if regDT.T != TK_REG {
				return tokens, fmt.Errorf("Destination should be R0-R15, found %v", regDT)
			}
			func5 = 11
		} else if regST.T == TK_REG_SP { //LD REG, SP
			if regDT.T != TK_REG {
				return tokens, fmt.Errorf("Destination should be R0-R15, found %v", regDT)
			}
		}
	} else { //Alu instruction
		if err != nil {
			return tokens, err
		}
		func5, err = getFunc5FromALU(tokens[0].Tk)
		if err != nil {
			return tokens, err
		}
	}
	instr.Func5 = func5
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.Instructions = append(parser.Instructions, instr)
	parser.CurAddress += 2
	return tokens, nil
}

func ParseLdRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, bankd, func2 uint8
	var err error

	ldT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	rd, bankd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	if ldT.T == TK_LD_1 {
		func2 = bankd<<1 | 1
	} else {
		func2 = bankd << 1
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd
	instr.Imm = int16(immT.ValInt)
	instr.Address = parser.CurAddress
	instr.Func2 = func2
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.Instructions = append(parser.Instructions, instr)
	parser.CurAddress += 2
	return tokens, nil
}

// ALU REG, IMM (7 bit Immediate operations)
//
//	func(2 bit) = register bank (0 or 1)
//	func(0-1 bits):
//
// 0) ADD/SUB: Rd = Rd + (signed(IMM))
// 1) SHL/SHR: Rd = << or >> signed(IMM & 31)
// 2) LDI Rd = [PC - IMM]; 32 bit constant loading, IMM in 32 bit dword (-512 ... 0 bytes)
// 3) DJNZ Rd, PC + signed(IMM); Rd-- if not zero, jump taken. IMM in instructions (-64 ... +63 instructions)
func ParseAluRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	var rd, bankd, func3 uint8
	var err error

	aluT := tokens[tokenpos]
	regDT := tokens[tokenpos+1]
	immT := tokens[tokenpos+3]
	instr := Instruction{}

	rd, bankd, err = parseRegister(regDT.Tk)
	if err != nil {
		return tokens, err
	}
	func3 = bankd << 2

	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
	} else {
		instr.Imm = int16(immT.ValInt)
	}

	switch aluT.Tk {
	case "ADD":
		func3 = 0
	case "SUB":
		instr.Imm = ^instr.Imm + 1
		func3 = 0
	case "SHL":
		func3 = 1
	case "SHR":
		instr.Imm = ^instr.Imm + 1
		func3 = 1
	case "LDI":
		func3 = 2
	case "DJNZ":
		func3 = 3
	default:
		return tokens, fmt.Errorf("Instruction should be one of: ADD, SUB, SHL, SHR, LDI, DJNZ! '%s' not found", aluT.Tk)
	}

	instr.Opcode = rule.Opcode
	instr.Rd = rd

	instr.Address = parser.CurAddress
	instr.Func3 = func3
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+4:]...)

	parser.Instructions = append(parser.Instructions, instr)
	parser.CurAddress += 2
	return tokens, nil
}

func ParseLabel(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	labelT := tokens[tokenpos]
	parser.Labels[labelT.Tk] = parser.CurAddress

	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, nil
}
