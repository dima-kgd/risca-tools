package isa

import (
	"fmt"
)

type Parser struct {
	Instructions []Instruction
	Memory       map[uint32]Instruction
	Labels       map[string]uint32
	Constants    map[string]int32
	CurAddress   uint32
}

func NewParser() *Parser {
	return &Parser{
		Instructions: make([]Instruction, 0, 1024),
		Memory:       make(map[uint32]Instruction),
		Labels:       make(map[string]uint32),
		Constants:    make(map[string]int32),
		CurAddress:   0,
	}
}

func (parser *Parser) ParseLine(line string) (bool, error) {
	tokens, err := Tokenize(line)
	if err != nil { //Error parsing tokens
		return false, fmt.Errorf("Format error: %v", err)
	}
	if len(tokens) == 0 { //Skip line
		return true, nil
	}

	var matched bool

	for len(tokens) > 0 {
		for _, rule := range syntaxRules {
			tokens, matched, err = parser.applyRule(rule, tokens)
			if err != nil {
				return false, fmt.Errorf("Syntax error: %v", err)
			}
			if matched {
				break
			}
		}
		if !matched {
			return false, fmt.Errorf("Syntax error: Unknown syntax: %s", line)
		}
	}
	return false, err
}

func tokenIn(token Token, tokens []uint8) bool {
	for _, t := range tokens {
		if token.T == t {
			return true
		}
	}
	return false
}

func (parser *Parser) addInstruction(instruction Instruction) error {
	if _, ok := parser.Memory[instruction.Address]; ok {
		return fmt.Errorf("Duplicate instruction at address: %d. Instructinon: %s", instruction.Address, instruction)
	}
	parser.Instructions = append(parser.Instructions, instruction)
	parser.Memory[instruction.Address] = instruction
	return nil
}

func (parser *Parser) applyRule(rule Rule, tokens []Token) ([]Token, bool, error) {
	pos := 0
	err := error(nil)
	for pos+len(rule.Syntax)-1 < len(tokens) {
		match := true
		for i, sxtokens := range rule.Syntax {
			if !tokenIn(tokens[pos+i], sxtokens) {
				match = false
				break
			}
		}
		if match {
			tokens, err = rule.ParseFunc(parser, rule, tokens, pos)
			return tokens, true, err
		} else {
			pos++
		}
	}
	return tokens, false, nil
}

func (parser *Parser) ProcessLabels() error {
	for i, instr := range parser.Instructions {
		if instr.Label != "" {
			addr, ok := parser.Labels[instr.Label]
			if !ok {
				return fmt.Errorf("Label not defined: %s", instr.Label)
			}
			offset := (int32(addr) - int32(instr.Address)) >> 1 //offset in instructions (2 bytes per instruction)
			parser.Instructions[i].Imm = int16(offset)
			switch instr.Opcode.Type {
			case OP_TYPE_7_IMM_REG:
				if f := instr.Func3 & 0b0000_0011; instr.Opcode.Opc == OP_ALU_REG_IMM && f > 1 { // LDI and DJNZ
					if f == 2 { //LDI
						if offset > 0 {
							return fmt.Errorf("Offset should be negative for DJNZ! %s, %s", instr.Label, instr)
						}
						parser.Instructions[i].Imm = int16(-offset)
					} else { //DJNZ - default offset
					}
				} else {
					if offset < -64 || offset > 63 {
						return fmt.Errorf("Label too large for instruction! %s, %s", instr.Label, instr)
					}
				}
			case OP_TYPE_8_IMM_REG:
				if offset < -128 || offset > 127 {
					return fmt.Errorf("Label too large for instruction! %s, %s", instr.Label, instr)
				}
			case OP_TYPE_13_IMM:
				if offset < -4096 || offset > 4095 {
					return fmt.Errorf("Label too large for instruction! %s, %s", instr.Label, instr)
				}
			case OP_TYPE_3_REG:
				if instr.Rd == 0 {
					if offset < -512 || offset > 511 {
						return fmt.Errorf("Label too large for instruction! %s, %s", instr.Label, instr)
					}
					parser.Instructions[i].Rs = uint8(offset >> 7)
					parser.Instructions[i].Ex = uint8(offset>>5) & 0b0000_0011
					parser.Instructions[i].Func2 = uint8(offset>>3) & 0b0000_0011
					parser.Instructions[i].Rx = uint8(offset) & 0b0000_0111
				} else {
					//TODO:
				}
			default:
				return fmt.Errorf("Label not applicable for instruction! %s, %s", instr.Label, instr)
			}
		}
	}
	return nil
}
