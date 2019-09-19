package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func RunParser(p Parser, input string) (interface{}, error) {
	state := NewParserState(input,0)
	return p(state).extract()
}

//找字符串的函数，这个还是自己写好一点
//func firstNonmatchingIndex(s1, s2 string, offset int) int {
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
//}

//这个函数不一定需要
//func DefaultSucceed(a interface{}) Parser {
	//return Str("").Map(func(i interface{}) interface{} {
	//	return a
	//})
//}

//直接获得一个成功的Parser
func Succeed(a interface{}) Parser {
	return func(state ParseState) ParseResult {
		return &Success{a,0}
	}
}

//读字符串
func Str(s string) Parser {
	return func(state ParseState) ParseResult {
		offset := state.loc.offset
		input := state.loc.input[offset:]
		if len(input) <= 0 || len(input) < len(s) || input == "\\z" {
			//文本已结束
			currentLine,line,col := getLineAndCol(state.loc.input,state.loc.offset)
			msg := "attach eof"
			err := NewParseError(currentLine,msg,line,col)
			return &Failure{err}
		}
		inputHead := input[:len(s)]
		if s == inputHead {
			//匹配到字符串
			return &Success{s,len(s)}
		} else {
			//匹配不到字符串
			currentLine,line,col := getLineAndCol(state.loc.input,state.loc.offset)
			msg := fmt.Sprintf("except:%s", s)
			err := NewParseError(currentLine,msg,line,col)
			return &Failure{err}
		}
	}
}

//读字符
func Char(c rune) Parser {
	s := string(c)
	return Str(s)
}

//读字符，存在的字符
func CharIn(cs []rune) Parser {
	return func(state ParseState) ParseResult {
		offset := state.loc.offset
		input := state.loc.input[offset:]
		inputChars := []rune(input)
		char := inputChars[0]
		for _, i := range cs {
			if i == char {
				//满足的字符
				s := string(char)
				return &Success{s, len(s)}
			}
 		}
		//匹配不到
		currentLine,line,col := getLineAndCol(state.loc.input,state.loc.offset)
		msg := fmt.Sprintf("except: chars in %s", string(cs))
		err := NewParseError(currentLine,msg,line,col)
		return &Failure{err}
	}
}

//读字符，不存在字符
func CharNotIn(cs []rune) Parser {
	return func(state ParseState) ParseResult {
		offset := state.loc.offset
		input := state.loc.input[offset:]
		inputChars := []rune(input)
		char := inputChars[0]
		for _, i := range cs {
			if i == char {
				//匹配不到
				currentLine,line,col := getLineAndCol(state.loc.input,state.loc.offset)
				msg := fmt.Sprintf("not except: chars in %s", string(cs))
				err := NewParseError(currentLine,msg,line,col)
				return &Failure{err}
			}
		}
		//满足的字符
		s := string(char)
		return &Success{s, len(s)}
	}
}

//读一个或多个的Parser
func Many1(p LazyParser) Parser {
	return func(state ParseState) ParseResult {
		res := p()(state)
		if res.isSuccess() {
			length := res.getLength()
			value := res.getSuccess()
			nextState := NewParserState(state.loc.input, state.loc.offset + length)
			//第一个解析成功时，就调用many1, 结果只能是success
			nextRes := Many(p)(nextState)
			sliceValue := nextRes.getSuccess()
			if list,ok := sliceValue.([]interface{}); ok {
				newList := []interface{}{value}
				newList = append(newList, list...)
				return &Success{newList, length + nextRes.getLength()}
			} else {
				panic("parse value wrong type")
			}
		} else {
			return res
		}
	}
}

//读零个或多个的Parser
func Many(p LazyParser) Parser {
	return func(state ParseState) ParseResult {
		res := p()(state)
		if res.isSuccess() {
			length := res.getLength()
			value := res.getSuccess()
			nextState := NewParserState(state.loc.input, state.loc.offset + length)
			//递归继续
			nextRes := Many(p)(nextState)
			if !nextRes.isSuccess() {
				//一下个不成功，就返回只有当前value的Success
				oneSlice := []interface{}{value}
				return &Success{oneSlice,length} //长度就是当前value长度
			} else {
				//下一个已经成功, Success{ slice, 那个成功的长度 }
				sucLen := nextRes.getLength()
				sucVal := nextRes.getSuccess()
				if list,ok := sucVal.([]interface{}); ok {
					newList := []interface{}{value}
					newList = append(newList, list...)
					return &Success{newList, length + sucLen}
				} else {
					panic("parse value wrong type")
				}
			}
		} else {
			//找不到就直接返回一个空的slice
			var emptySlice []interface{}
			return &Success{emptySlice,0}
		}
	}
}

//正则读
func Regex(s string) Parser {
	msg := fmt.Sprintf("regex '%s'", s)
	r, _ := regexp.Compile(s)
	return func(state ParseState) ParseResult {
		input := state.loc.input[state.loc.offset:]
		if res, ok := r.FindString(input), r.MatchString(input); ok {
			if res == input[:len(res)] {
				return &Success{res, len(res)}
			}
		}
		currentLine, line, col := getLineAndCol(state.loc.input, state.loc.offset)
		msg := fmt.Sprintf("except:%s from %s", s, msg)
		err := NewParseError(currentLine, msg, line, col)
		return &Failure{err}
	}
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
	return Str("\"").
		SkipL( MakeLazyParser(Thru("\"")) ).
		Map(func(i interface{}) interface{} {
		if str, ok := i.(string) ;ok {
			chars := []rune(str)
			chars = chars[:len(chars)-1]
			str = string(chars)
			return str
		} else {
			panic("parse value wrong type")
		}
	})
}

//读转义字符
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

//读带转义字符的字符串（忽略头尾引号)
func EscapedQuoted() Parser {
	return Surround(
		MakeLazyParser(Str("\"")),
		MakeLazyParser(Str("\"")),
		MakeLazyParser(Many1(
			MakeLazyParser(Escaped().Or(
				MakeLazyParser(CharNotIn([]rune("\"\\")))))))).
		Map(func(i interface{}) interface{} {
			// Many(Escaped.Or(CharNotIn(..))) ,Escaped得到的是string，CharNotIn得到的也是string，所以这里得到的是 []string
			if is,ok := i.([]interface{}); ok {
				var ss []string
				for _, s := range is {
					if str,ok2 := s.(string); ok2 {
						ss = append(ss,str)
					} else {
						panic("parse value wrong type")
					}
				}
				res := strings.Join(ss,"")
				return res
			} else {
				panic("parse value wrong type")
			}
	})
}

//读浮点型字符串
func DoubleString() Parser {
	return Token(
		MakeLazyParser(
			Regex("[-+]?([0-9]*\\.)?[0-9]+([eE][-+]?[0-9]+)?")))
}

//读出浮点数
func Double() Parser {
	return DoubleString().Map(func(i interface{}) interface{} {
		if str,ok := i.(string); ok {
			double,err := strconv.ParseFloat(str, 64)
			if err != nil {
				panic(err)
			}
			return double
		} else {
			panic("parse value wrong type")
		}
	})
}

//token是前后带空格的（或者没有空格也行）
func Token(p LazyParser) Parser {
	return Surround(Whitespace,Whitespace,p)
}

//前后被包围
func Surround(start, stop ,p LazyParser) Parser {
	return start().SkipL(p).SkipR(stop)
}

//读取到结束符
func Eof() Parser {
	return Regex("\\z").Label("unexpected trailing characters")
}

//默认的根，省略最后的结束符
func Root(p Parser) Parser {
	return p.SkipR(Eof)
}
