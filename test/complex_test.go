package test

import (
	"SimpleGoParser/parser"
	"testing"
)

func TestMany(t *testing.T) {
	input := "some text"
	p := parser.CharIn([]rune{'s','o','m','e'})
	p2 := parser.Many(parser.MakeLazyParser(p))
	t.Log(RunParser(p2,input))
	input2 := "abcd some text"
	t.Log(RunParser(p2,input2))
}

func TestMany1(t *testing.T) {
	input := "some text"
	p := parser.CharIn([]rune{'s','o','m','e'})
	p2 := parser.Many1(parser.MakeLazyParser(p))
	t.Log(RunParser(p2,input))
	input2 := "abcd some text"
	t.Log(RunParser(p2,input2))
}

func TestQuote(t *testing.T) {
	input := "\"some text\""
	p := parser.Quoted()
	t.Log(RunParser(p,input))
}

func TestEscaped(t *testing.T) {
	input := "\\\"" +
		"\\\\" +
		"\\/" +
		"\\b" +
		"\\f" +
		"\\n" +
		"\\r" +
		"\\t"
	p := parser.Escaped()
	p2 := parser.Many1(parser.MakeLazyParser(p))
	t.Log(RunParser(p2,input))
	input2 := "\\\"\\t"
	t.Log(RunParser(p2,input2))
}

func TestEscapedQuoted(t *testing.T) {
	input := "\"this \\\"is\\\" \\t name \""
	p := parser.EscapedQuoted()
	t.Log(RunParser(p,input))
}

func TestToken(t *testing.T) {
	input := "  some text   "
	p := parser.Regex("[A-z]*")
	p2 := parser.Token(parser.MakeLazyParser(p))
	t.Log(RunParser(p2,input))
}

func TestSurround(t *testing.T) {
	start := parser.Str("{")
	p := parser.Regex("[A-z]*")
	p2 := parser.Token(parser.MakeLazyParser(p))
	end := parser.Str("}")
	surround := parser.Surround(parser.MakeLazyParser(start),parser.MakeLazyParser(end),parser.MakeLazyParser(p2))
	input := "{ some  test }"
	input2 := "{ some   }"
	t.Log(RunParser(surround,input))
	t.Log(RunParser(surround,input2))
}