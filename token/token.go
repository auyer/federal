// Copyright (c) 2014, Rob Thornton
// All rights reserved.
// This source code is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package token

type Token int

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

var tok_strings = map[Token]string{
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

func (t Token) IsLiteral() bool {
	return t > litStart && t < litEnd
}

func (t Token) IsOperator() bool {
	return t > opStart && t < opEnd
}

func (t Token) IsKeyword() bool {
	return t > keyStart && t < keyEnd
}

func Lookup(str string) Token {
	if str == "true" || str == "false" {
		return BOOL
	}
	if str == "add" {
		return ADD
	}
	if str == "mul" {
		return MUL
	}
	if str == "div" {
		return DIV
	}
	if str == "sub" {
		return SUB
	}
	if str == "rem" {
		return REM
	}
	if str == "compare" {
		return EQLLIT
	}
	if str == "do" {
		return DO
	}
	for t, s := range tok_strings {
		if s == str {
			return t
		}
	}
	return IDENT
}

func (t Token) String() string {
	return tok_strings[t]
}

func (t Token) Valid() bool {
	return t > tokenStart && t < tokenEnd
}
