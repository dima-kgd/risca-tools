package isa

import (
	"fmt"
	"strconv"
	"strings"
)

func Tokenize(input string) ([]Token, error) { // input line in uppercase, e.g. "ADD R1, R2"
	input = strings.ToUpper(input)
	var tokens []Token
	pos := 0
	for pos < len(input) {
		if isWhitespace(input, pos) {
			pos = skipWhitespace(input, pos)
		} else if isWord(input, pos) {
			word, newPos := readWord(input, pos)
			pos = newPos
			if tokenType, exists := mapWordToToken[word]; exists {
				tokens = append(tokens, Token{T: tokenType, Tk: word})
			} else {
				tokens = append(tokens, Token{T: TK_LABEL, Tk: word})
				//return nil, fmt.Errorf("unknown token: %s", word)
			}
		} else if isNumber(input, pos) {
			number, newPos, isParsed := readNumber(input, pos)
			if !isParsed {
				return nil, fmt.Errorf("Not a number on pos %d: %s>%s", pos, input[:pos], input[pos:])
			}
			tokens = append(tokens, Token{T: TK_NUMBER, Tk: input[pos:newPos], ValInt: number})
			pos = newPos
		} else {
			var token Token
			var found bool
			token, pos, found = getComma(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getPlusMinus(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			token, pos, found = getBracket(input, pos)
			if found {
				tokens = append(tokens, token)
				continue
			}
			return nil, fmt.Errorf("Unexpeced symbol on pos %d: %c. %s>%s", pos, input[pos], input[:pos], input[pos:])
		}
	}
	return tokens, nil
}

func getPlusMinus(input string, pos int) (Token, int, bool) {
	switch input[pos] {
	case '+':
		pos++
		return Token{T: TK_PLUS, Tk: "+"}, pos, true
	case '-':
		pos++
		return Token{T: TK_MINUS, Tk: "-"}, pos, true
	default:
		return Token{}, pos, false
	}

}

func getBracket(input string, pos int) (Token, int, bool) {
	switch input[pos] {
	case '[':
		pos++
		return Token{T: TK_L_SQBR, Tk: "["}, pos, true
	case ']':
		pos++
		return Token{T: TK_R_SQBR, Tk: "]"}, pos, true
	default:
		return Token{}, pos, false
	}
}

func getComma(input string, pos int) (Token, int, bool) {
	if input[pos] == ',' {
		pos++
		return Token{T: TK_COMMA, Tk: ","}, pos, true
	}
	return Token{}, pos, false
}

func isWord(input string, pos int) bool {
	return (input[pos] >= 'A' && input[pos] <= 'Z') || (input[pos] == '.')
}

func isNumber(input string, pos int) bool {
	if isDigit(input, pos) {
		return true
	}
	if input[pos] == '0' && len(input) > pos+1 && input[pos+1] == 'X' {
		return true
	}
	return false
}

func isDigit(input string, pos int) bool {
	return input[pos] >= '0' && input[pos] <= '9'
}

func isWhitespace(input string, pos int) bool {
	return input[pos] == ' ' || input[pos] == '\t' || input[pos] == ';'
}

func skipWhitespace(input string, pos int) int {
	for pos < len(input) && (input[pos] == ' ' || input[pos] == '\t' || input[pos] == ';') {
		if input[pos] == ';' { // Skip line comment
			return len(input)
		} else {
			pos++
		}
	}
	return pos
}

func readWord(input string, pos int) (string, int) { // Read word started from letter, can contain letters, digits and dots (R1, R2, etc.)
	start := pos
	for pos < len(input) && (isWord(input, pos) || isDigit(input, pos) || input[pos] == '.') {
		pos++
	}
	return input[start:pos], pos
}

func isHexDigit(input string, pos int) bool {
	return (input[pos] >= '0' && input[pos] <= '9') || (input[pos] >= 'A' && input[pos] <= 'F')
}

func readNumber(input string, pos int) (int, int, bool) { // Read number started from digit
	if input[pos] == '0' && len(input) > pos+1 && input[pos+1] == 'X' { //Read hex
		pos += 2 // Skip '0X' prefix
		start := pos
		for pos < len(input) && isHexDigit(input, pos) {
			pos++
		}
		value, _ := strconv.ParseInt(input[start:pos], 16, 64)
		return int(value), pos, true
	} else { //Read decimal
		start := pos
		for pos < len(input) && isDigit(input, pos) {
			pos++
		}
		value, err := strconv.Atoi(input[start:pos])
		if err != nil {
			return 0, pos, false
		}
		return value, pos, true
	}

}
