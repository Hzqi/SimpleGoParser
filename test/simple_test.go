package test

import (
	"SimpleGoParser/parser"
	"testing"
)

//基本测试，不含符合的Parser的

func RunParser(p parser.Parser, input string) (interface{}, error) {
	return parser.RunParser(p,input)
}

func TestSucceed(t *testing.T) {
	input := "some text"
	p := parser.Succeed("a")
	t.Log(RunParser(p,input))
}

func TestStr(t *testing.T) {
	input := "some text"
	p := parser.Str("some")
	t.Log(RunParser(p,input))
	p2 := parser.Str("abc")
	t.Log(RunParser(p2,input))
}

func TestRegex(t *testing.T) {
	input := "some text"
	p := parser.Regex("[A-z]*\\ [A-z]")
	t.Log(RunParser(p,input))
	input2 := "123some text"
	p2 := parser.Regex("[A-z]*\\ [A-z]")
	t.Log(RunParser(p2,input2))
}

func TestWhitespace(t *testing.T) {
	input := "some text"
	p := parser.Whitespace()
	t.Log(RunParser(p,input))
	input2 := "   some text"
	p2 := parser.Whitespace()
	t.Log(RunParser(p2,input2))
}

func TestDigits(t *testing.T) {
	input := "some text"
	p := parser.Digits()
	t.Log(RunParser(p,input))
	input2 := "12345 are digits"
	t.Log(RunParser(p,input2))
}

func TestThur(t *testing.T) {
	input := "\"some\" text"
	p := parser.Thru("\"")
	t.Log(RunParser(p,input))
	p2 := parser.Thru("x")
	t.Log(RunParser(p2,input))
}

func TestDoubleString(t *testing.T) {
	input := "123.456 text"
	p := parser.DoubleString()
	t.Log(RunParser(p,input))
	input2 := "1.234e5"
	t.Log(RunParser(p,input2))
}

func TestDouble(t *testing.T)  {
	input := "123.456 text"
	p := parser.Double()
	t.Log(RunParser(p,input))
	input2 := "1.234e5"
	t.Log(RunParser(p,input2))
}

func TestChar(t *testing.T) {
	input := "some text"
	p := parser.Char('s')
	p2 := parser.Char('a')
	t.Log(RunParser(p,input))
	t.Log(RunParser(p2,input))
}

func TestCharIn(t *testing.T) {
	input := "some text"
	p := parser.CharIn([]rune{'s','a'})
	p2 := parser.CharIn([]rune{'b','a'})
	t.Log(RunParser(p,input))
	t.Log(RunParser(p2,input))
}

func TestCharNotIn(t *testing.T) {
	input := "some text"
	p := parser.CharNotIn([]rune{'s','a'})
	p2 := parser.CharNotIn([]rune{'b','a'})
	t.Log(RunParser(p,input))
	t.Log(RunParser(p2,input))
}