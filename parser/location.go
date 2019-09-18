package parser

import "strings"

type Location struct {
	input  string
	offset int
}

func NewLocation(input string, offset int) *Location {
	return &Location{input, offset}
}

func NewLocationDefault(input string) *Location {
	return NewLocation(input, 0)
}

//记录第几行
func (loc *Location) line() int {
	current := loc.input[0 : loc.offset+1]
	var line = 0
	for _, ch := range []rune(current) {
		if ch == '\n' {
			line += 1
		}
	}
	return line
}

//记录第几个
func (loc *Location) col() int {
	current := loc.input[0 : loc.offset+1]
	lstIdsNL := strings.LastIndex(current, "\n")
	if lstIdsNL == -1 {
		return loc.offset + 1
	} else {
		return loc.offset - lstIdsNL
	}
}

func (loc *Location) toError(msg string) ParseError {
	return ParseError{[]LocStr{LocStr{*loc, msg}}}
}

func (loc *Location) advanceBy(n int) *Location {
	newloc := Location{
		input:  loc.input,
		offset: loc.offset + n,
	}
	return &newloc
}

func (loc *Location) currentLine() string {
	if len(loc.input) > 1 {
		lines := strings.Split(loc.input, "\n")
		return lines[loc.line()]
	} else {
		return ""
	}
}

func (loc *Location) columnCaret() string {
	return strings.Repeat(" ", loc.col()-1) + "^"
}
