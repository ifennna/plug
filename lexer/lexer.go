package lexer

import "plug/token"

type Lexer struct {
	input           string
	currentPosition int
	readPosition    int
	currentChar     byte
}

func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	switch lexer.currentChar {
	case '=':
		tok = newToken(token.ASSIGN, lexer.currentChar)
	case '+':
		tok = newToken(token.PLUS, lexer.currentChar)
	case '(':
		tok = newToken(token.LPAREN, lexer.currentChar)
	case ')':
		tok = newToken(token.RPAREN, lexer.currentChar)
	case '{':
		tok = newToken(token.LBRACE, lexer.currentChar)
	case '}':
		tok = newToken(token.RBRACE, lexer.currentChar)
	case ',':
		tok = newToken(token.COMMA, lexer.currentChar)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.currentChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.currentChar = 0
	} else {
		lexer.currentChar = lexer.input[lexer.readPosition]
	}

	lexer.currentPosition = lexer.readPosition
	lexer.readPosition++
}

func newToken(tokenType token.Type, character byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(character)}
}
