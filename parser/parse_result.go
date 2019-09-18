package parser

type ParseResult interface {
	isSuccess() bool
	//success
	getSuccess() interface{}
	getLength() int
	//failure
	getError() ParseError
}

type Success struct {
	value  interface{}
	length int
}

type Failure struct {
	parseError ParseError
	committed  bool
}

