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

func firstNonmatchingIndex(s1, s2 string, offset int) int {
	i := 0
	for i < len(s1) &&
		i < len(s2) &&
		(i+offset) < len(s1) {
		if s1[i+offset] != s2[i] {
			return i
		} else {
			i += 1
		}
	}
	if len(s1)-offset >= len(s2) {
		return -1
	} else {
		return len(s1) - offset
	}
}

func DefaultSucceed(a interface{}) Parser {
	return Str("").Map(func(i interface{}) interface{} {
		return a
	})
}

func Succeed(a interface{}) Parser {
	return func(state ParseState) ParseResult {
		return &Success{a, 0}
	}
}

func Str(s string) Parser {
	msg := fmt.Sprintf("'%s'", s)
	return func(state ParseState) ParseResult {
		idx := firstNonmatchingIndex(state.loc.input, s, state.loc.offset)
		if idx == -1 {
			//找到了
			return &Success{s, len(s)}
		} else {
			return &Failure{state.loc.advanceBy(idx).toError(msg), idx != 0}
		}
	}
}

func Char(c rune) Parser {
	return Str(string(c)).Map(func(i interface{}) interface{} {
		if str, ok := i.(string); ok {
			return []rune(str)[0]
		} else {
			panic("parser result wrong type")
		}
	})
}

func CharIn(cs []rune) Parser {
	return func(state ParseState) ParseResult {
		offset := state.loc.offset
		char := []rune(state.loc.input[offset:])[0]
		contains := func(list []rune, e rune) bool {
			for _, v := range list {
				if v == e {
					return true
				}
			}
			return false
		}
		if contains(cs, char) {
			return &Success{char, 1}
		} else {
			return &Failure{state.loc.advanceBy(state.loc.col()).toError("except: " + string(char)), false}
		}
	}
}

func CharNotIn(cs []rune) Parser {
	return func(state ParseState) ParseResult {
		offset := state.loc.offset
		char := []rune(state.loc.input[offset:])[0]
		contains := func(list []rune, e rune) bool {
			for _, v := range list {
				if v == e {
					return true
				}
			}
			return false
		}
		if !contains(cs, char) {
			return &Success{char, 1}
		} else {
			return &Failure{state.loc.advanceBy(state.loc.col()).toError("not except: " + string(char)), false}
		}
	}
}

func Many1(p LazyParser) LazyParser {
	return func() Parser {
		return p().Map2(Many(p), func(a interface{}, b interface{}) interface{} {
			if list, ok := a.([]interface{}); ok {
				return append(list, b)
			} else {
				panic("parser result wrong type")
			}
		})
	}
}

func Many(p LazyParser) LazyParser {
	lazySucceed := func() Parser{
		return Succeed(make([]interface{}, 5))
	}
	return func() Parser {
		return p().Map2(Many(p), func(a interface{}, b interface{}) interface{} {
			if list, ok := a.([]interface{}); ok {
				return append(list, b)
			} else {
				panic("parser result wrong type")
			}
		}).Or(lazySucceed)
	}
}

func Regex(s string) Parser {
	msg := fmt.Sprintf("regex '%s'", s)
	r, _ := regexp.Compile(s)
	return func(state ParseState) ParseResult {
		if res, ok := r.FindString(state.input()), r.MatchString(state.input()); ok {
			return &Success{res, len(res)}
		} else {
			return &Failure{state.loc.toError(msg), false}
		}
	}
}

func Whitespace() Parser {
	return Regex("\\s*")
}

func Digits() Parser {
	return Regex("\\d+")
}

func Thru(s string) Parser {
	return Regex(".*?" + regexp.QuoteMeta(s))
}

func Quoted() Parser {
	lazy := func() Parser{
		return Thru("\"").Map(func(i interface{}) interface{} {
			if str, ok := i.(string); ok {
				rs := []rune(str)
				rs = rs[:len(rs)-1]
				return string(rs)
			} else {
				panic("parser result wrong type")
			}
		})
	}
	return Str("\"").SkipL(lazy)
}

func Escaped() Parser {
	becomefunc := func(str string) func(interface{}) interface{} {
		return func(i interface{}) interface{} {
			return str
		}
	}
	a := Str("\\\"").Map(becomefunc("\""))
	b := func() Parser{return Str("\\\\").Map(becomefunc("\\"))}
	c := func() Parser{return Str("\\/").Map(becomefunc("/"))}
	d := func() Parser{return Str("\\b").Map(becomefunc("\b"))}
	e := func() Parser{return Str("\\f").Map(becomefunc("\f"))}
	f := func() Parser{return Str("\\n").Map(becomefunc("\n"))}
	g := func() Parser{return Str("\\r").Map(becomefunc("\r"))}
	h := func() Parser{return Str("\\t").Map(becomefunc("\t"))}
	return a.Or(b).Or(c).Or(d).Or(e).Or(f).Or(g).Or(h)
}

func EscapedQuoted() Parser {
	lazy := func() Parser{
		return CharNotIn([]rune("\"\\"))
	}
	p := Many1(func() Parser {
		return Escaped().Or(lazy)
	})
	return Surround(Str("\""), Str("\""), p).Map(func(i interface{}) interface{} {
		if list, ok := i.([]interface{}); ok {
			ss := make([]rune, len(list))
			for i, v := range list {
				if s, ok2 := v.(rune); ok2 {
					ss[i] = s
				} else {
					panic("parser result wrong type")
				}
			}
			return string(ss)
		} else {
			panic("parser result wrong type")
		}
	})
}

func DoubleString() Parser {
	return Token(Regex("[-+]?([0-9]*\\.)?[0-9]+([eE][-+]?[0-9]+)?"))
}

func Double() Parser {
	return DoubleString().Map(func(i interface{}) interface{} {
		if str, ok := i.(string); ok {
			f, err := strconv.ParseFloat(str, 64)
			if err != nil {
				panic(err)
			}
			return f
		} else {
			panic("parser result wrong type")
		}
	}).Label("double literal")
}

func Token(p Parser) Parser {
	lazy := func() Parser{
		return Whitespace()
	}
	lazyp := func() Parser{
		return p
	}
	return Whitespace().SkipL(lazyp).SkipR(lazy)
}

func Surround(start, stop Parser, p LazyParser) Parser {
	lazy := func() Parser{
		return stop
	}
	return start.SkipL(p).SkipR(lazy)
}

func Eof() Parser {
	return Regex("\\z").Label("unexpected trailing characters")
}

func Root(p Parser) Parser {
	return p.SkipR(Eof)
}
