package pos

import (
	"fmt"
	"path/filepath"
)

type PosType byte

const (
	InFilePos PosType = iota
	InBuiltInPos
	InConsolePos
)

type Pos struct {
	line, column uint
	path         string
	t            PosType
	relative     *struct{ line, column uint }
}

func BuiltInPos(f string) Pos {
	return Pos{
		path: f,
		t:    InBuiltInPos,
	}
}
func ConsolePos(line, column uint) Pos {
	return Pos{
		line:   line,
		column: column,
		t:      InConsolePos,
	}
}
func NewRelative(relative Pos, line, column uint) Pos {
	relative.relative = &struct {
		line   uint
		column uint
	}{line, column}
	return relative
}
func NewPos(line, column uint, path string) (Pos, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Pos{}, err
	}
	return Pos{line, column, absPath, InFilePos, nil}, nil
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
	if p.relative != nil {
		p.relative.column += n
		return
	}
	p.column += n
}
func (p *Pos) AddOneColumn() {
	if p.relative != nil {
		p.relative.column++
		return
	}
	p.column++
}
func (p *Pos) NextLine() {
	if p.relative != nil {
		p.relative.line++
		p.relative.column = 0
		return
	}
	p.line++
	p.column = 0
}

func (p Pos) String() string {
	if p.relative != nil {
		line, column := p.relative.line, p.relative.column
		p.relative = nil
		return fmt.Sprintf("%s (position in the string: %d:%d)", p.String(), line, column)
	}
	if p.t == InConsolePos {
		return fmt.Sprintf("%d:%d", p.line, p.column)
	}
	if p.t == InBuiltInPos {
		return fmt.Sprintf("build-in function %s", p.path)
	}
	return fmt.Sprintf("%s:%d:%d", p.path, p.line, p.column)
}
