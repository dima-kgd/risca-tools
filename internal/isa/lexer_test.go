package isa

import (
	"reflect"
	"testing"
)

func TestLexer_Tokenize(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "ldi reg label",
			input:    "ldi r0 variable1",
			expected: []Token{{T: TK_LDI, Tk: "LDI"}, {T: TK_REG, Tk: "R0"}, {T: TK_LABEL, Tk: "VARIABLE1"}},
		},
		{
			name:     "djnz reg label",
			input:    "djnz r15 loop1",
			expected: []Token{{T: TK_DJNZ, Tk: "DJNZ"}, {T: TK_REG, Tk: "R15"}, {T: TK_LABEL, Tk: "LOOP1"}},
		},
		{
			name:     "djnz reg imm",
			input:    "djnz r15 -20",
			expected: []Token{{T: TK_DJNZ, Tk: "DJNZ"}, {T: TK_REG, Tk: "R15"}, {T: TK_MINUS, Tk: "-"}, {T: TK_NUMBER, Tk: "20", ValInt: 20}},
		},
		{
			name:     "jmp reg cmp reg, -imm",
			input:    "jmp r5 EQ r6, -16",
			expected: []Token{{T: TK_JMP, Tk: "JMP"}, {T: TK_REG, Tk: "R5"}, {T: TK_EQ, Tk: "EQ"}, {T: TK_REG, Tk: "R6"}, {T: TK_COMMA, Tk: ","}, {T: TK_MINUS, Tk: "-"}, {T: TK_NUMBER, Tk: "16", ValInt: 16}},
		},
		{
			name:     "jmp reg cmp reg, reg",
			input:    "jmp r5 NE r6, r10",
			expected: []Token{{T: TK_JMP, Tk: "JMP"}, {T: TK_REG, Tk: "R5"}, {T: TK_NE, Tk: "NE"}, {T: TK_REG, Tk: "R6"}, {T: TK_COMMA, Tk: ","}, {T: TK_REG, Tk: "R10"}},
		},
		{
			name:     "ret",
			input:    "\tret\t",
			expected: []Token{{T: TK_RET, Tk: "RET"}},
		},
		{
			name:     "call reg",
			input:    "call r7",
			expected: []Token{{T: TK_CALL, Tk: "CALL"}, {T: TK_REG, Tk: "R7"}},
		},
		{
			name:     "jmp reg",
			input:    "jmp r4",
			expected: []Token{{T: TK_JMP, Tk: "JMP"}, {T: TK_REG, Tk: "R4"}},
		},
		{
			name:     "ld.b reg, [reg+imm]",
			input:    "ld.b r10, [r12+5]",
			expected: []Token{{T: TK_LD_BYTE, Tk: "LD.B"}, {T: TK_REG, Tk: "R10"}, {T: TK_COMMA, Tk: ","}, {T: TK_L_SQBR, Tk: "["}, {T: TK_REG, Tk: "R12"}, {T: TK_PLUS, Tk: "+"}, {T: TK_NUMBER, Tk: "5", ValInt: 5}, {T: TK_R_SQBR, Tk: "]"}},
		},
		{
			name:     "ld.w reg, [reg++]",
			input:    "ld.w r10, [r12++]",
			expected: []Token{{T: TK_LD_WORD, Tk: "LD.W"}, {T: TK_REG, Tk: "R10"}, {T: TK_COMMA, Tk: ","}, {T: TK_L_SQBR, Tk: "["}, {T: TK_REG, Tk: "R12"}, {T: TK_PLUS, Tk: "+"}, {T: TK_PLUS, Tk: "+"}, {T: TK_R_SQBR, Tk: "]"}},
		},
		{
			name:     "ld.b reg, imm",
			input:    "ld.1 r6, 10",
			expected: []Token{{T: TK_LD_1, Tk: "LD.1"}, {T: TK_REG, Tk: "R6"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "10", ValInt: 10}},
		},
		{
			name:     "add reg, imm",
			input:    "add r6, 1",
			expected: []Token{{T: TK_ALU, Tk: "ADD"}, {T: TK_REG, Tk: "R6"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "1", ValInt: 1}},
		},
		{
			name:     "ld.b reg, [reg]",
			input:    "ld.b r11, [r12]",
			expected: []Token{{T: TK_LD_BYTE, Tk: "LD.B"}, {T: TK_REG, Tk: "R11"}, {T: TK_COMMA, Tk: ","}, {T: TK_L_SQBR, Tk: "["}, {T: TK_REG, Tk: "R12"}, {T: TK_R_SQBR, Tk: "]"}},
		},
		{
			name:     "pop lr",
			input:    "pop lr",
			expected: []Token{{T: TK_POP, Tk: "POP"}, {T: TK_REG_LR, Tk: "LR"}},
		},
		{
			name:     "push reg",
			input:    "push r11",
			expected: []Token{{T: TK_PUSH, Tk: "PUSH"}, {T: TK_REG, Tk: "R11"}},
		},
		{
			name:     "add reg reg",
			input:    "\tADD R1, R2",
			expected: []Token{{T: TK_ALU, Tk: "ADD"}, {T: TK_REG, Tk: "R1"}, {T: TK_COMMA, Tk: ","}, {T: TK_REG, Tk: "R2"}},
		},
		{
			name:     "ld reg reg",
			input:    "\tLD R9, R0",
			expected: []Token{{T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R9"}, {T: TK_COMMA, Tk: ","}, {T: TK_REG, Tk: "R0"}},
		},
		{
			name:     "ld reg imm dec",
			input:    "\tLD R3, 255",
			expected: []Token{{T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R3"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "255", ValInt: 255}},
		},
		{
			name:     "ld reg imm hex",
			input:    "\tLD R3, 0xFF",
			expected: []Token{{T: TK_LD, Tk: "LD"}, {T: TK_REG, Tk: "R3"}, {T: TK_COMMA, Tk: ","}, {T: TK_NUMBER, Tk: "0XFF", ValInt: 255}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := Tokenize(tc.input)
			if err != nil {
				t.Errorf("Tokenize(%q) error: %v", tc.input, err)
			}
			if !reflect.DeepEqual(tokens, tc.expected) {
				t.Errorf("Tokenize(%q) = %v; want %v", tc.input, tokens, tc.expected)
			}
		})
	}
}
