// package token is responsible for storing all tokens and basic information on recognizing them.

package token

type Token int

// all constants recieve a integer from the iota declaration.
// This means it's possible to add new keyworkds without having to change the corresponding numbers.
const (
	tokenStart Token = iota

	EOF
	ILLEGAL
	COMMENT

	litStart
	BOOL
	IDENT
	INTEGER
	litEnd

	opStart
	LKEY
	RKEY
	LPAREN
	RPAREN
	ADD
	ADDLIT
	SUB
	SUBLIT
	MUL
	MULLIT
	DIV
	DIVLIT
	REM
	REMLIT
	EQL
	EQLLIT
	ASSIGN
	opEnd

	keyStart
	DO
	FOR
	FUNC
	IF
	VAR
	OR
	keyEnd

	tokenEnd
)

// tokStrings stores the correspondance between keywords and tokens
var tokStrings = map[Token]string{
	EOF:     "EOF",
	ILLEGAL: "Illegal",
	COMMENT: "Comment",
	INTEGER: "Integer",
	LKEY:    "{",
	RKEY:    "}",
	LPAREN:  "(",
	RPAREN:  ")",
	ADD:     "add",
	ADDLIT:  "+",
	SUB:     "sub",
	SUBLIT:  "-",
	MUL:     "mul",
	MULLIT:  "*",
	DIV:     "div",
	DIVLIT:  "/",
	DO:      "do",
	REM:     "rem",
	REMLIT:  "%",
	EQL:     "compare",
	EQLLIT:  "==",
	ASSIGN:  "=",
}

// checks if a token is inside the literal block
func (t Token) IsLiteral() bool {
	return t > litStart && t < litEnd
}

// checks if a token is inside the operator block
func (t Token) IsOperator() bool {
	return t > opStart && t < opEnd
}

// checks if a token is inside the keyword block
func (t Token) IsKeyword() bool {
	return t > keyStart && t < keyEnd
}

// checks if a token is the String tokens list
func Lookup(str string) Token {
	if str == "true" || str == "false" {
		return BOOL
	}
	// if str == "add" {
	// 	return ADD
	// }
	// if str == "mul" {
	// 	return MUL
	// }
	// if str == "div" {
	// 	return DIV
	// }
	// if str == "sub" {
	// 	return SUB
	// }
	// if str == "rem" {
	// 	return REM
	// }
	// if str == "compare" {
	// 	return EQLLIT
	// }
	// if str == "do" {
	// 	return DO
	// }
	for t, s := range tokStrings {
		if s == str {
			return t
		}
	}
	return IDENT
}

// String function allows printing tokens and basic conversion
func (t Token) String() string {
	return tokStrings[t]
}

// checks if a token is inside the token block
func (t Token) Valid() bool {
	return t > tokenStart && t < tokenEnd
}
