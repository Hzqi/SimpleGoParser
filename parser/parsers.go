package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

func RunParser(p Parser, input string) Either {
	state := NewParseState(*NewLocationDefault(input))
	return extract(p(state))
}

//找字符串的函数，这个还是自己写好一点
func firstNonmatchingIndex(s1, s2 string, offset int) int {
	//i := 0
	//for i < len(s1) &&
	//	i < len(s2) &&
	//	(i+offset) < len(s1) {
	//	if s1[i+offset] != s2[i] {
	//		return i
	//	} else {
	//		i += 1
	//	}
	//}
	//if len(s1)-offset >= len(s2) {
	//	return -1
	//} else {
	//	return len(s1) - offset
	//}
}

//这个函数不一定需要
func DefaultSucceed(a interface{}) Parser {
	//return Str("").Map(func(i interface{}) interface{} {
	//	return a
	//})
}

//直接获得一个成功的Parser
func Succeed(a interface{}) Parser {

}

//读字符串
func Str(s string) Parser {

}

//读字符
func Char(c rune) Parser {

}

//读字符，存在的字符
func CharIn(cs []rune) Parser {

}

//读字符，不存在字符
func CharNotIn(cs []rune) Parser {

}

//读一个或多个的Parser
func Many1(p LazyParser) Parser {

}

//读零个或多个的Parser
func Many(p LazyParser) Parser {

}

//正则读
func Regex(s string) Parser {
	//msg := fmt.Sprintf("regex '%s'", s)
	//r, _ := regexp.Compile(s)
	//return func(state ParseState) ParseResult {
	//	if res, ok := r.FindString(state.input()), r.MatchString(state.input()); ok {
	//		return &Success{res, len(res)}
	//	} else {
	//		return &Failure{state.loc.toError(msg), false}
	//	}
	//}
}

//读若干个连续的空格
func Whitespace() Parser {
	return Regex("\\s*")
}

//读数字
func Digits() Parser {
	return Regex("\\d+")
}

//一直读取字符串，直到遇到指定的字符串为止
func Thru(s string) Parser {
	return Regex(".*?" + regexp.QuoteMeta(s))
}

//读带引号的字符串（忽略头尾引号）
func Quoted() Parser {

}

//读转义字符
func Escaped() Parser {
	//becomefunc := func(str string) func(interface{}) interface{} {
	//	return func(i interface{}) interface{} {
	//		return str
	//	}
	//}
	//a := Str("\\\"").Map(becomefunc("\""))
	//b := func() Parser{return Str("\\\\").Map(becomefunc("\\"))}
	//c := func() Parser{return Str("\\/").Map(becomefunc("/"))}
	//d := func() Parser{return Str("\\b").Map(becomefunc("\b"))}
	//e := func() Parser{return Str("\\f").Map(becomefunc("\f"))}
	//f := func() Parser{return Str("\\n").Map(becomefunc("\n"))}
	//g := func() Parser{return Str("\\r").Map(becomefunc("\r"))}
	//h := func() Parser{return Str("\\t").Map(becomefunc("\t"))}
	//return a.Or(b).Or(c).Or(d).Or(e).Or(f).Or(g).Or(h)
}

//读带转义字符的字符串（忽略头尾引号)
func EscapedQuoted() Parser {

}

//读浮点型字符串
func DoubleString() Parser {
	return Token(Regex("[-+]?([0-9]*\\.)?[0-9]+([eE][-+]?[0-9]+)?"))
}

//读出浮点数
func Double() Parser {

}

//token是前后带空格的（或者没有空格也行）
func Token(p Parser) Parser {

}

//前后被包围
func Surround(start, stop Parser, p LazyParser) Parser {
	lazy := func() Parser{
		return stop
	}
	return start.SkipL(p).SkipR(lazy)
}

//读取到结束符
func Eof() Parser {
	return Regex("\\z").Label("unexpected trailing characters")
}

//默认的根，省略最后的结束符
func Root(p Parser) Parser {
	return p.SkipR(Eof)
}
