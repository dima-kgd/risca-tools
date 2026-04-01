package isa

import "fmt"

const (
	RuleALURegReg = iota
	LdRegImm
	AluRegImm
	Label
	Equ
	EvalConst
	DefineDBVar
	JmpRel
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
var constEquSyntax = [][]uint8{{TK_LABEL}, {TK_EQU}, {TK_NUMBER}}
var evalConstSyntax = [][]uint8{{TK_BUCKS}, {TK_LABEL}}
var defineVarSyntax = [][]uint8{{TK_LABEL}, {TK_DB}, {TK_NUMBER, TK_STRING}}
var jmpRelSyntax = [][]uint8{{TK_JMP}, {TK_NUMBER, TK_LABEL}}

var syntaxRules = []Rule{
	{Type: Equ, Syntax: constEquSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseConstEqu},
	{Type: EvalConst, Syntax: evalConstSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseEvalConst},
	{Type: Label, Syntax: labelSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseLabel},
	{Type: DefineDBVar, Syntax: defineVarSyntax, Opcode: GetOpcode(OP_DB), ParseFunc: parseDefineDbVar},
	{Type: RuleALURegReg, Syntax: aluRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG), ParseFunc: parseAluRegReg},
	{Type: LdRegImm, Syntax: ldRegImmSyntax, Opcode: GetOpcode(OP_LD_REG_IMM), ParseFunc: parseLdRegImm},
	{Type: AluRegImm, Syntax: aluRegImmSyntax, Opcode: GetOpcode(OP_ALU_REG_IMM), ParseFunc: parseAluRegImm},
	{Type: JmpRel, Syntax: jmpRelSyntax, Opcode: GetOpcode(OP_JUMP_REL), ParseFunc: parseJmpRel},
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

func parseAluRegReg(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
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

	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

func parseLdRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
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

	parser.addInstruction(instr)
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
// 2) 2) LDI Rd = PC - IMM; get 32 bit address relative to PC, IMM in instructions (-128 ... 0 instructions)
// 3) DJNZ Rd, PC + signed(IMM); Rd-- if not zero, jump taken. IMM in instructions (-64 ... +63 instructions)
func parseAluRegImm(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
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

	parser.addInstruction(instr)
	parser.CurAddress += 2
	return tokens, nil
}

func parseJmpRel(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	immT := tokens[tokenpos+1]
	instr := Instruction{}
	instr.Rd = 0
	if immT.T == TK_LABEL {
		instr.Label = immT.Tk
	} else {
		instr.Rs = uint8(immT.ValInt >> 7)
		instr.Ex = uint8(immT.ValInt>>5) & 0b0000_0011
		instr.Func2 = uint8(immT.ValInt>>3) & 0b0000_0011
		instr.Rx = uint8(immT.ValInt) & 0b0000_0111
	}
	instr.Opcode = rule.Opcode
	instr.Address = parser.CurAddress
	parser.addInstruction(instr)
	parser.CurAddress += 2

	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, nil
}

func addLabel(parser *Parser, token Token) error {
	if _, ok := parser.Labels[token.Tk]; ok {
		return fmt.Errorf("Label '%s' already defined! On address: 0x%08x", token.Tk, parser.Labels[token.Tk])
	}
	parser.Labels[token.Tk] = parser.CurAddress
	return nil
}

//Example: label1:
func parseLabel(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	if err := addLabel(parser, tokens[tokenpos]); err != nil {
		return tokens, err
	}
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+2:]...)
	return tokens, nil
}

//Example: myconst1 equ 0x1234
func parseConstEqu(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	constT := tokens[tokenpos]
	numberT := tokens[tokenpos+2]
	if _, ok := parser.Constants[constT.Tk]; ok {
		return tokens, fmt.Errorf("Constant '%s' already defined! Value: 0x%08x", constT.Tk, parser.Constants[constT.Tk])
	}
	parser.Constants[constT.Tk] = numberT.ValInt
	//Remove parsed tokens from the list
	tokens = append(tokens[:tokenpos], tokens[tokenpos+3:]...)
	return tokens, nil
}

// Example: add r1, $const1; replace $const1 by number
func parseEvalConst(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	constT := tokens[tokenpos+1]
	constVal, ok := parser.Constants[constT.Tk]
	if !ok {
		return tokens, fmt.Errorf("Constant '%s' is not defined!", constT.Tk)
	}
	//Remove second token and replace first by number
	tokens[tokenpos] = Token{T: TK_NUMBER, Tk: string(constVal), ValInt: constVal}
	tokens = append(tokens[:tokenpos+1], tokens[tokenpos+2:]...)
	return tokens, nil
}

//Example: Mystr db 'Hello',0
func parseDefineDbVar(parser *Parser, rule Rule, tokens []Token, tokenpos int) ([]Token, error) {
	labelT := tokens[tokenpos]
	instr := Instruction{}
	instrByte := 1

	if err := addLabel(parser, labelT); err != nil {
		return tokens, err
	}

	for i := tokenpos + 2; i < len(tokens); i++ {
		numberT := tokens[i]
		switch numberT.T {
		case TK_COMMA:
			continue
		case TK_NUMBER:
			if numberT.ValInt < 0 || numberT.ValInt > 255 {
				return tokens, fmt.Errorf("Number out of range(0..255) for db: 0x%08X", numberT.ValInt)
			}
			if instrByte == 1 {
				instr.Opcode = rule.Opcode
				instr.Address = parser.CurAddress
				instr.Imm = instr.Imm | int16(numberT.ValInt)<<8
				instrByte = 0
			} else {
				instr.Imm = instr.Imm | int16(numberT.ValInt)
				instrByte = 1
				parser.addInstruction(instr)
				parser.CurAddress += 2
				instr = Instruction{}
			}
		case TK_STRING:
			for _, char := range []byte(numberT.ValStr) {
				if instrByte == 1 {
					instr.Opcode = rule.Opcode
					instr.Address = parser.CurAddress
					instr.Imm = instr.Imm | int16(char)<<8
					instrByte = 0
				} else {
					instr.Imm = instr.Imm | int16(char)
					instrByte = 1
					parser.addInstruction(instr)
					parser.CurAddress += 2
					instr = Instruction{}
				}
			}
		default:
			return tokens, fmt.Errorf("Unexpected token: %s", numberT)
		}
	}
	if instrByte == 0 {
		parser.addInstruction(instr)
		parser.CurAddress += 2
	}
	tokens = []Token{}
	return tokens, nil
}
