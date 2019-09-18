package parser

type ParseError struct {
	stack []LocStr
}

type LocStr struct {
	loc Location
	str string
}