package scan

import (
	"fmt"

	"github.com/auyer/federal/token"
)

type Error struct {
	pos token.Position
	msg string
}

func (e Error) Error() string {
	return fmt.Sprint(e.pos, " ", e.msg)
}

type ErrorList []*Error

func (el ErrorList) Count() int {
	return len(el)
}

func (el *ErrorList) Add(p token.Position, msg string) {
	*el = append(*el, &Error{p, msg})
}

func (el *ErrorList) Print() {
	for _, err := range *el {
		fmt.Println(err)
	}
}
