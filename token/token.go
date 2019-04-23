package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	SLASH    = "/"
	ASTERISK = "*"

	LT = "<"
	GT = ">"

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]Type{
	"func": FUNCTION,
	"let":  LET,
}

func LookUpIdentifier(identifier string) Type {
	if token, exists := keywords[identifier]; exists {
		return token
	}

	return IDENTIFIER
}
