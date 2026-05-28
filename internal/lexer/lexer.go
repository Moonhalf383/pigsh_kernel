package lexer

import (
	"strings"
	"unicode"
)

// TokenKind represents the type of a lexical token.
type TokenKind int

const (
	TokenKeyword   TokenKind = iota // var, add, print, etc.
	TokenIdent                      // variable or label name
	TokenIntLit                     // integer literal
	TokenFloatLit                   // float literal
	TokenStringLit                  // "hello" or 'hello'
	TokenBoolLit                    // true, false
)

// Token is a single lexical token with its kind and literal text.
type Token struct {
	Kind    TokenKind
	Literal string
	Line    int // 1-based source line number
}

// Tokenize splits pigsh source code into a stream of tokens.
// Each source line becomes one instruction's worth of tokens.
// Empty lines and comment lines (# ...) are skipped.
func Tokenize(source string) []Token {
	lines := strings.Split(source, "\n")
	var tokens []Token

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lineTokens := tokenizeLine(line, lineNum+1)
		tokens = append(tokens, lineTokens...)
	}

	return tokens
}

func tokenizeLine(line string, lineNum int) []Token {
	var tokens []Token
	i := 0
	runes := []rune(line)

	for i < len(runes) {
		ch := runes[i]

		// skip whitespace
		if unicode.IsSpace(ch) {
			i++
			continue
		}

		// string literal (double or single quoted)
		if ch == '"' || ch == '\'' {
			tok, newI := scanString(runes, i, lineNum)
			tokens = append(tokens, tok)
			i = newI
			continue
		}

		// word (keyword, ident, number, bool)
		if !unicode.IsSpace(ch) {
			tok, newI := scanWord(runes, i, lineNum)
			tokens = append(tokens, tok)
			i = newI
			continue
		}

		i++
	}

	return tokens
}

func scanString(runes []rune, start int, lineNum int) (Token, int) {
	quote := runes[start]
	i := start + 1
	var buf []rune

	for i < len(runes) && runes[i] != quote {
		buf = append(buf, runes[i])
		i++
	}

	if i < len(runes) {
		i++ // skip closing quote
	}

	return Token{Kind: TokenStringLit, Literal: string(buf), Line: lineNum}, i
}

func scanWord(runes []rune, start int, lineNum int) (Token, int) {
	i := start
	var buf []rune

	for i < len(runes) && !unicode.IsSpace(runes[i]) {
		buf = append(buf, runes[i])
		i++
	}

	word := string(buf)

	// check bool literals
	if word == "true" || word == "false" {
		return Token{Kind: TokenBoolLit, Literal: word, Line: lineNum}, i
	}

	// check number
	if kind, ok := numberKind(word); ok {
		return Token{Kind: kind, Literal: word, Line: lineNum}, i
	}

	// keywords are the instruction mnemonics
	if isKeyword(word) {
		return Token{Kind: TokenKeyword, Literal: word, Line: lineNum}, i
	}

	// otherwise it's an identifier (variable or label name)
	return Token{Kind: TokenIdent, Literal: word, Line: lineNum}, i
}

func numberKind(s string) (TokenKind, bool) {
	if s == "" {
		return 0, false
	}

	i := 0
	if s[0] == '-' || s[0] == '+' {
		if len(s) == 1 {
			return 0, false
		}
		i = 1
	}

	hasDot := false
	for ; i < len(s); i++ {
		if s[i] == '.' {
			if hasDot {
				return 0, false
			}
			hasDot = true
		} else if s[i] < '0' || s[i] > '9' {
			return 0, false
		}
	}

	if hasDot {
		return TokenFloatLit, true
	}
	return TokenIntLit, true
}

var keywords = map[string]bool{
	"var": true, "add": true, "sub": true, "mul": true, "div": true,
	"mod": true, "and": true, "or": true, "xor": true, "not": true,
	"mov": true, "beq": true, "bne": true, "blt": true, "bgt": true,
	"ble": true, "bge": true, "jump": true, "label": true, "print": true,
	"input": true, "push": true, "pop": true, "call": true, "ret": true,
	"inc": true, "dec": true, "neg": true, "halt": true, "nop": true,
}

func isKeyword(s string) bool {
	return keywords[s]
}
