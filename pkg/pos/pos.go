package pos

import (
	"fmt"
	"path/filepath"
)

type Pos struct {
	line, column uint
	path         string
	buildIn      bool
}

func BuildInPos(f string) Pos {
	return Pos{
		path:    f,
		buildIn: true,
	}
}

func NewPos(line, column uint, path string) (Pos, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Pos{}, err
	}
	return Pos{line, column, absPath, false}, nil
}
func (p Pos) GetLine() uint {
	return p.line
}
func (p Pos) GetColumn() uint {
	return p.column
}
func (p Pos) GetPath() string {
	return p.path
}

func (p *Pos) AddColumn(n uint) {
	p.column += n
}
func (p *Pos) AddOneColumn() {
	p.column++
}
func (p *Pos) NextLine() {
	p.line++
	p.column = 0
}

func (p Pos) String() string {
	if p.buildIn {
		return fmt.Sprintf("build-in function %s", p.path)
	}
	return fmt.Sprintf("%s:%d:%d", p.path, p.line, p.column)
}
