package isa

import "fmt"

type Token struct {
	Type        uint8
	TokenString string
	ValueInt    uint
	ValueStr    string
}

func (token Token) String() string {
	return fmt.Sprintf("%s(%s)", GetTokenTypeString(token.Type), token.TokenString)
}

const (
	TK_LD           = 0x00
	TK_ST           = 0x01
	TK_LD_DOT_BYTE  = 0x02
	TK_LD_DOT_WORD  = 0x03
	TK_LD_DOT_DWORD = 0x04
	TK_REG          = 0x05
	TK_REG_SP       = 0x06
	TK_REG_LR       = 0x07
	TK_COMMA        = 0x08
	TK_ALU          = 0x09
	TK_NUMBER       = 0x0A
)

var tokenTypeToString = map[uint8]string{
	TK_LD:           "TK_LD",
	TK_ST:           "TK_ST",
	TK_LD_DOT_BYTE:  "TK_DOT_BYTE",
	TK_LD_DOT_WORD:  "TK_DOT_WORD",
	TK_LD_DOT_DWORD: "TK_DOT_DWORD",
	TK_REG:          "TK_REG",
	TK_REG_SP:       "TK_REG_SP",
	TK_REG_LR:       "TK_REG_LR",
	TK_COMMA:        "TK_COMMA",
	TK_ALU:          "TK_ALU",
	TK_NUMBER:       "TK_NUMBER",
}

func GetTokenTypeString(tokenType uint8) string {
	if tokenString, exists := tokenTypeToString[tokenType]; exists {
		return tokenString
	}
	return "TK_UNKNOWN"
}

var mapRegisterToNumber = map[string]uint8{
	"R0":  0,
	"R1":  1,
	"R2":  2,
	"R3":  3,
	"R4":  4,
	"R5":  5,
	"R6":  6,
	"R7":  7,
	"R8":  8,
	"R9":  9,
	"R10": 10,
	"R11": 11,
	"R12": 12,
	"R13": 13,
	"R14": 14,
	"R15": 15,
}

func GetRegisterNumber(reg string) (uint8, bool) {
	if regNum, exists := mapRegisterToNumber[reg]; exists {
		return regNum, true
	}
	return 0, false
}

var mapWordToToken = map[string]uint8{
	"LD":   TK_LD,
	"ST":   TK_ST,
	"LD.B": TK_LD_DOT_BYTE,
	"LD.W": TK_LD_DOT_WORD,
	"LD.D": TK_LD_DOT_DWORD,
	"R0":   TK_REG,
	"R1":   TK_REG,
	"R2":   TK_REG,
	"R3":   TK_REG,
	"R4":   TK_REG,
	"R5":   TK_REG,
	"R6":   TK_REG,
	"R7":   TK_REG,
	"R8":   TK_REG,
	"R9":   TK_REG,
	"R10":  TK_REG,
	"R11":  TK_REG,
	"R12":  TK_REG,
	"R13":  TK_REG,
	"R14":  TK_REG,
	"R15":  TK_REG,
	"SP":   TK_REG_SP,
	"LR":   TK_REG_LR,
	",":    TK_COMMA,
	"ADD":  TK_ALU,
	"SUB":  TK_ALU,
	"SHL":  TK_ALU,
	"SHR":  TK_ALU,
	"AND":  TK_ALU,
	"OR":   TK_ALU,
	"XOR":  TK_ALU,
	"NOT":  TK_ALU,
	"MUL":  TK_ALU,
	"INT":  TK_ALU,
}

func GetTokenTypeByWord(word string) (uint8, bool) {
	if tokenType, exists := mapWordToToken[word]; exists {
		return tokenType, true
	}
	return 0, false
}
