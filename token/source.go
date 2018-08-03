package token

import "fmt"

type Source struct {
	base  int
	name  string
	src   string
	lines []int
}

func NewSource(name, src string) *Source {
	return &Source{
		base:  1,
		name:  name,
		src:   src,
		lines: make([]int, 0, 16),
	}
}

func (f *Source) AddLine(offset int) {
	if offset >= f.base-1 && offset < f.base+len(f.src) {
		f.lines = append(f.lines, offset)
	}
}

func (f *Source) Base() int {
	return f.base
}

func (f *Source) Pos(offset int) Pos {
	if offset < 0 || offset >= len(f.src) {
		panic("illegal file offset")
	}
	return Pos(f.base + offset)
}

func (f *Source) Position(p Pos) Position {
	col, row := int(p), 1

	for i, nl := range f.lines {
		if p > f.Pos(nl) {
			col, row = int(p-f.Pos(nl)), i+1
		}
	}

	return Position{Filename: f.name, Col: col, Row: row}
}

func (f *Source) Size() int {
	return len(f.src)
}

type Pos uint

var illegalPos = Pos(0)

func (p Pos) Valid() bool {
	return p != illegalPos
}

type Position struct {
	Filename string
	Col, Row int
}

func (p Position) String() string {
	if p.Filename == "" {
		return fmt.Sprintf("%d:%d", p.Row, p.Col)
	}
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Row, p.Col)
}
