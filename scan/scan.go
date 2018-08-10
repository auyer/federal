// package scan is responsible for retrieving all runes from the text in the source file
package scan

import (
	"unicode"

	"github.com/auyer/federal/token"
)

// Scanner structure stores all data necessary for the scanner
type Scanner struct {
	ch      rune
	offset  int
	roffset int
	src     string
	file    *token.Source
}

func (s *Scanner) Init(file *token.Source, src string) {
	s.file = file
	s.offset, s.roffset = 0, 0
	s.src = src
	s.file.AddLine(s.offset)

	s.next()
}

// Scan will analyse a rune and return its features
func (s *Scanner) Scan() (lit string, tok token.Token, pos token.Pos) {
	s.skipWhitespace()

	if unicode.IsLetter(s.ch) {
		return s.scanIdentifier()
	}

	if unicode.IsDigit(s.ch) {
		return s.scanNumber()
	}

	lit, pos = string(s.ch), s.file.Pos(s.offset)
	switch s.ch {
	case '{':
		tok = token.LKEY
	case '}':
		tok = token.RKEY
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case '+':
		tok = token.ADDLIT
	case '-':
		tok = token.SUBLIT
	case '*':
		tok = token.MULLIT
	case '=':
		tok = s.selectToken('=', token.EQL, token.ASSIGN)
	case '/':
		tok = s.selectToken('/', token.COMMENT, token.DIV)
		if tok == token.COMMENT {
			s.skipComment()
			s.next()
			return s.Scan()
		}
	case '%':
		tok = token.REMLIT
	default:
		if s.offset >= len(s.src)-1 {
			tok = token.EOF
		} else {
			tok = token.ILLEGAL
		}
	}

	s.next()

	return
}

// next will move all cursors to the next rune
func (s *Scanner) next() {
	s.ch = rune(0)
	if s.roffset < len(s.src) {
		s.offset = s.roffset
		s.ch = rune(s.src[s.offset])
		if s.ch == '\n' {
			s.file.AddLine(s.offset)
		}
		s.roffset++
	}
}

// scanIdentifier will create a string from runes
func (s *Scanner) scanIdentifier() (string, token.Token, token.Pos) {
	start := s.offset

	for unicode.IsLetter(s.ch) || unicode.IsDigit(s.ch) {
		s.next()
	}
	offset := s.offset
	if s.ch == rune(0) {
		offset++
	}
	lit := s.src[start:offset]
	return lit, token.Lookup(lit), s.file.Pos(start)
}

// scanIdentifier will create a string from numerical runes
func (s *Scanner) scanNumber() (string, token.Token, token.Pos) {
	start := s.offset

	for unicode.IsDigit(s.ch) {
		s.next()
	}
	offset := s.offset
	if s.ch == rune(0) {
		offset++
	}
	return s.src[start:offset], token.INTEGER, s.file.Pos(start)
}

// skipComment will move the cursors to the first non-comment rune
func (s *Scanner) skipComment() {
	for s.ch != '\n' && s.offset < len(s.src)-1 {
		s.next()
	}
}

// skipWhitespace will move the cursors to the first non-space rune
func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.ch) {
		s.next()
	}
}

// readNext reads the next token without changing the state of nything else. Use to indentify Multi rune literals
func (s *Scanner) readNext() rune {
	if s.roffset < len(s.src) {
		return rune(s.src[s.roffset])
	}
	return rune(0)
}

// selectToken checks if the next rune is == to r, and if so, sets the TOKEN to a. Otherwise, to b
func (s *Scanner) selectToken(r rune, a, b token.Token) token.Token {
	nrune := s.readNext()
	if nrune == r {
		s.next()
		return a
	}
	return b
}
