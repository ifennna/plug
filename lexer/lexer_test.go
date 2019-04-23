package lexer

import (
	"fmt"
	"plug/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for index, testToken := range tests {
		resultToken := lexer.NextToken()

		if resultToken.Type != testToken.expectedType {
			t.Fatalf("tests[%d] - tokentype is wrong, expected=%q, got=%q",
				index, testToken.expectedType, resultToken.Type)
		}

		if resultToken.Literal != testToken.expectedLiteral {
			t.Fatalf("tests[%d] - literal is wrong, expected=%q, got=%q",
				index, testToken.expectedLiteral, resultToken.Literal)
		}

		fmt.Print(resultToken)
	}
}
