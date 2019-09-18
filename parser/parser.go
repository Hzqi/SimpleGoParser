package parser

//实际上是一个 ParseState => ParseResult 的函数
type Parser func(state ParseState) ParseResult

type LazyParser func()Parser

func MakeLazyParser(p Parser) func()Parser {
	return func() Parser {
		return p
	}
}