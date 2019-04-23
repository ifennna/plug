package lexer

import (
	"plug/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `let five = 5;
let ten = 10;

let add = func(x, y) {
	x + y;
};

let result = add(five, ten);
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
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
	}
}
