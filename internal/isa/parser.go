package isa

import (
	"fmt"
)

func ParseLine(line string) (Instruction, error) {
	tokens, err := Tokenize(line)
	if err != nil {
		return Instruction{}, err
	}
	if len(tokens) == 0 {
		return Instruction{}, fmt.Errorf("Token expected")
	}

	var expectedToken uint8
	var errorToken uint8
	var matched bool
	var matchedRule Rule
	for _, rule := range syntaxRules {
		if tokens[0].Type == rule.Syntax[0] {
			matched, expectedToken, errorToken = ruleMatchesTokens(rule, tokens)
			if matched {
				matchedRule = rule
				break
			}
		}
	}
	if !matched {
		return Instruction{}, fmt.Errorf("Syntax error: expected %s, got %s", GetTokenTypeString(expectedToken), GetTokenTypeString(errorToken))
	}
	return parseInstruction(matchedRule, tokens)
}

func ruleMatchesTokens(rule Rule, tokens []Token) (bool, uint8, uint8) {
	for i, tokenType := range rule.Syntax {
		if tokens[i].Type != tokenType {
			return false, tokenType, tokens[i].Type
		}
	}
	return true, 0, 0
}

func parseRegisters(tokenrd, tokenrs string) (uint8, uint8, uint8, uint8, error) {
	bankd, banks := uint8(0), uint8(0)
	rd, found := GetRegisterNumber(tokenrd)
	if !found {
		return 0, 0, 0, 0, fmt.Errorf("Invalid register: %s", tokenrd)
	}
	rs, found := GetRegisterNumber(tokenrs)
	if !found {
		return 0, 0, 0, 0, fmt.Errorf("Invalid register: %s", tokenrs)
	}
	if rd > 7 {
		rd = rd - 7
		bankd = 1
	}
	if rs > 7 {
		rs = rs - 7
		banks = 1
	}
	return rd, bankd, rs, banks, nil
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

func makeEx(bankd, banks uint8) uint8 {
	return (banks << 1) | bankd
}

func makeAluRegRegInstruction(rule Rule, tokens []Token) (Instruction, error) {
	var rd, bankd, rs, banks uint8
	var err error
	var func5 uint8
	switch rule.Type {
	case RuleALURegReg:
		rd, bankd, rs, banks, err = parseRegisters(tokens[1].TokenString, tokens[3].TokenString)
		if err != nil {
			return Instruction{}, err
		}
		func5, err = getFunc5FromALU(tokens[0].TokenString)
	case RuleLDRegReg:
		rd, bankd, rs, banks, err = parseRegisters(tokens[1].TokenString, tokens[3].TokenString)
		func5 = 0
	case RuleLDRegSP:
		rd, bankd, err = parseRegister(tokens[1].TokenString)
		func5 = 10
	case RuleLDRegLR:
		rd, bankd, err = parseRegister(tokens[1].TokenString)
		func5 = 11
	case RuleLDSPReg:
		rs, banks, err = parseRegister(tokens[3].TokenString)
		func5 = 12
	case RuleLDLRReg:
		rs, banks, err = parseRegister(tokens[3].TokenString)
		func5 = 13
	}
	if err != nil {
		return Instruction{}, err
	}
	return Instruction{Opcode: rule.Opcode, Rd: rd, Rs: rs, Ex: makeEx(bankd, banks), Func5: func5}, nil
}

func parseInstruction(rule Rule, tokens []Token) (Instruction, error) {
	switch rule.Type {
	case RuleALURegReg, RuleLDRegReg, RuleLDRegSP, RuleLDRegLR, RuleLDSPReg, RuleLDLRReg:
		return makeAluRegRegInstruction(rule, tokens)
	default:
		return Instruction{}, fmt.Errorf("Unknown rule type: %d", rule.Type)
	}
}

const (
	RuleALURegReg = iota
	RuleLDRegReg
	RuleLDRegSP
	RuleLDRegLR
	RuleLDSPReg
	RuleLDLRReg
)

type Rule struct {
	Type   uint8
	Syntax []uint8
	Opcode Opcode
}

var aluRegRegSyntax = []uint8{TK_ALU, TK_REG, TK_COMMA, TK_REG}
var ldRegRegSyntax = []uint8{TK_LD, TK_REG, TK_COMMA, TK_REG}
var ldRegSPSyntax = []uint8{TK_LD, TK_REG, TK_COMMA, TK_REG_SP}
var ldRegLRSyntax = []uint8{TK_LD, TK_REG, TK_COMMA, TK_REG_LR}
var ldSPRegSyntax = []uint8{TK_LD, TK_REG_SP, TK_COMMA, TK_REG}
var ldLRRegSyntax = []uint8{TK_LD, TK_REG_LR, TK_COMMA, TK_REG}

var syntaxRules = []Rule{
	{Type: RuleALURegReg, Syntax: aluRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
	{Type: RuleLDRegReg, Syntax: ldRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
	{Type: RuleLDRegSP, Syntax: ldRegSPSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
	{Type: RuleLDRegLR, Syntax: ldRegLRSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
	{Type: RuleLDSPReg, Syntax: ldSPRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
	{Type: RuleLDLRReg, Syntax: ldLRRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
}
