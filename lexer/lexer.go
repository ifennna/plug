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

	lexer.skipWhitespace()

	switch lexer.currentChar {
	case '=':
		if lexer.peekChar() == '=' {
			character := lexer.currentChar
			lexer.readChar()
			literal := string(character) + string(lexer.currentChar)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, lexer.currentChar)
		}
	case '+':
		tok = newToken(token.PLUS, lexer.currentChar)
	case '-':
		tok = newToken(token.MINUS, lexer.currentChar)
	case '!':
		if lexer.peekChar() == '=' {
			character := lexer.currentChar
			lexer.readChar()
			literal := string(character) + string(lexer.currentChar)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, lexer.currentChar)
		}
	case '/':
		tok = newToken(token.SLASH, lexer.currentChar)
	case '*':
		tok = newToken(token.ASTERISK, lexer.currentChar)
	case '<':
		tok = newToken(token.LT, lexer.currentChar)
	case '>':
		tok = newToken(token.GT, lexer.currentChar)
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
	case '"':
		tok.Type = token.STRING
		tok.Literal = lexer.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lexer.currentChar) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookUpIdentifier(tok.Literal)
			return tok
		} else if isDigit(lexer.currentChar) {
			tok.Type = token.INT
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.currentChar)
		}
	}

	lexer.readChar()
	return tok
}

func (lexer *Lexer) skipWhitespace() {
	for lexer.currentChar == ' ' || lexer.currentChar == '\t' || lexer.currentChar == '\n' || lexer.currentChar == '\r' {
		lexer.readChar()
	}
}

func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
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

func isLetter(character byte) bool {
	return 'a' <= character && character <= 'z' || 'A' <= character && character <= 'Z' || character == '_'
}

func isDigit(character byte) bool {
	return '0' <= character && character <= '9'
}

func (lexer *Lexer) readIdentifier() string {
	position := lexer.currentPosition
	for isLetter(lexer.currentChar) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.currentPosition]
}

func (lexer *Lexer) readNumber() string {
	position := lexer.currentPosition
	for isDigit(lexer.currentChar) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.currentPosition]
}

func (lexer *Lexer) readString() string {
	position := lexer.currentPosition + 1

	for {
		lexer.readChar()
		if lexer.currentChar == '"' || lexer.currentChar == 0 {
			break
		}
	}

	return lexer.input[position:lexer.currentPosition]
}
