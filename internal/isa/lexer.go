package isa

import (
	"fmt"
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
				tokens = append(tokens, Token{Type: tokenType, TokenString: word})
			} else {
				return nil, fmt.Errorf("unknown token: %s", word)
			}
		} else if isComma(input, pos) {
			tokens = append(tokens, Token{Type: TK_COMMA, TokenString: ","})
			pos++
		} else {
			return nil, fmt.Errorf("Unexpeced symbol on pos %d: %c. %s>%s", pos, input[pos], input[:pos], input[pos:])
		}
	}
	return tokens, nil
}

func isWord(input string, pos int) bool {
	return (input[pos] >= 'A' && input[pos] <= 'Z') || (input[pos] == '.')
}

func isDigit(input string, pos int) bool {
	return input[pos] >= '0' && input[pos] <= '9'
}

func isWhitespace(input string, pos int) bool {
	return input[pos] == ' ' || input[pos] == '\t'
}

func isComma(input string, pos int) bool {
	return input[pos] == ','
}

func skipWhitespace(input string, pos int) int {
	for pos < len(input) && (input[pos] == ' ' || input[pos] == '\t') {
		pos++
	}

	return pos
}

func readWord(input string, pos int) (string, int) { // Read word started from letter, can contain letters and digits (R1, R2, etc.)
	start := pos
	for pos < len(input) && (isWord(input, pos) || isDigit(input, pos)) {
		pos++
	}
	return input[start:pos], pos
}
