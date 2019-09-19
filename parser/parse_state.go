package parser

//状态信息，其实就是Location
type ParseState struct {
	loc Location //每次的location都要是新的一个location
}

func NewParserState(input string, offset int) ParseState {
	loc := NewLocation(input,offset)
	return ParseState{loc}
}