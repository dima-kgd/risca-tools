package isa

import (
	"fmt"
	"strings"
)

func ParseLine(line string) (Instruction, error) {
	tokens, err := Tokenize(strings.ToUpper(line))
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
		return Instruction{}, fmt.Errorf("Syntax error: expected %s, got %s", GetTokenTypeString(expectedToken), errorToken)
	}
	return ParseInstruction(matchedRule, tokens)
}

func ruleMatchesTokens(rule Rule, tokens []Token) (bool, uint8, uint8) {
	for i, tokenType := range rule.Syntax {
		if tokens[i].Type != tokenType {
			return false, tokenType, tokens[i].Type
		}
	}
	return true, 0, 0
}

func ParseRegister(token string) (uint8, uint8, error) {
	reg, found := GetRegisterNumber(token)
	if !found {
		return 0, 0, fmt.Errorf("Invalid register: %s", token)
	}
	if reg > 7 {
		return reg - 7, 1, nil
	} else {
		return reg, 0, nil
	}
}

func ParseInstruction(rule Rule, tokens []Token) (Instruction, error) {
	switch rule.Type {
	case RULE_ALU_REG_REG:
		rd, bankd, err := ParseRegister(tokens[1].TokenString)
		if err != nil {
			return Instruction{}, err
		}
		rs, banks, err := ParseRegister(tokens[3].TokenString)
		if err != nil {
			return Instruction{}, err
		}
		return Instruction{Opcode: rule.Opcode, Rd: rd, Rs: rs, Ex: bankd<<1 | banks}, nil
	case RULE_LD_REG_REG:
		rd, found := GetRegisterNumber(tokens[1].TokenString)
		if !found {
			return Instruction{}, fmt.Errorf("Invalid register: %s", tokens[1].TokenString)
		}
		rs, found := GetRegisterNumber(tokens[3].TokenString)
		if !found {
			return Instruction{}, fmt.Errorf("Invalid register: %s", tokens[3].TokenString)
		}
		return Instruction{Opcode: rule.Opcode, Rd: rd, Rs: rs, Func5: 0}, nil
	default:
		return Instruction{}, fmt.Errorf("Unknown rule type: %d", rule.Type)
	}
}

const (
	RULE_ALU_REG_REG = 0
	RULE_LD_REG_REG  = 1
)

type Rule struct {
	Type   uint8
	Syntax []uint8
	Opcode Opcode
}

var aluRegRegSyntax = []uint8{TK_ALU, TK_REG, TK_COMMA, TK_REG}
var ldRegRegSyntax = []uint8{TK_LD, TK_REG, TK_COMMA, TK_REG}

var syntaxRules = []Rule{
	{Type: RULE_ALU_REG_REG, Syntax: aluRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
	{Type: RULE_LD_REG_REG, Syntax: ldRegRegSyntax, Opcode: GetOpcode(OP_ALU_REG_REG)},
}
