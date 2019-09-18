package parser

type ParseState struct {
	loc Location
}

func NewParseState(loc Location) ParseState {
	return ParseState{loc}
}

func (ps *ParseState) advance(numChars int) ParseState {
	loc := ps.loc
	loc.offset += numChars
	return ParseState{loc}
}

func (ps *ParseState) input() string {
	return ps.loc.input[ps.loc.offset:]
}

func (ps *ParseState) slice(n int) string {
	return ps.loc.input[ps.loc.offset : ps.loc.offset+n]
}
