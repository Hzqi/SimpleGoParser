package parser

//实际上是一个 ParseState => ParseResult 的函数
type Parser func(state ParseState) ParseResult

type LazyParser func()Parser

type mapFunc func(interface{})interface{}

type mapFunc2 func(interface{}, interface{}) interface{}

func MakeLazyParser(p Parser) func()Parser {
	return func() Parser {
		return p
	}
}

//Parser[A]
//带label的parser
func (p Parser) Label(name string) Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if res.isSuccess() {
			return res
		} else {
			err := res.getError()
			err.msg = name
			return res
		}
	}
}

//Parser[A]
//带scope的parser
func (p Parser) Scope(name string) Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if res.isSuccess() {
			return res
		} else {
			err := res.getError()
			err.msg = "in scope:" + name + " " + err.msg
			return res
		}
	}
}

//Parser[String]
//分割成功的字符串
func (p Parser) Slice() Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if res.isSuccess(){
			// if v,ok := res.getSuccess().(string); ok {
			// 	newStr := v[res.getLength():]
			// 	return &Success{newStr,res.getLength()}
			// } else {
			// 	panic("parse value wrong type")
			// }
			length := res.getLength()
			newStr := state.loc.input[state.loc.offset : state.loc.offset + length]
			return &Success{newStr, length}
		} else {
			return res
		}
	}
}

//Parser[B] A->B
//映射内容的Parser
func (p Parser) Map(f mapFunc) Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if res.isSuccess() {
			length := res.getLength()
			value := f(res.getSuccess())
			return &Success{value,length}
		} else {
			return res
		}
	}
}

//Parser[C] A,B -> C
//两变量的映射内容的Parser
func (p Parser) Map2(p2 LazyParser,f mapFunc2) Parser {
	return func(state ParseState) ParseResult {
		res1 := p(state)
		if res1.isSuccess() {
			afterRes := state.loc.offset + res1.getLength()
			state2 := NewParserState(state.loc.input, afterRes)
			res2 := p2()(state2)

			if res2.isSuccess() {
				length := res2.getLength() + res1.getLength()
				value := f(res1.getSuccess(), res2.getSuccess())
				return &Success{value,length}
			} else {
				return res2
			}
		} else {
			return res1
		}
	}
}

//Parser[(A,B)]
func (p Parser) Product(p2 LazyParser) Parser {
	return p.Map2(p2, func(a interface{}, b interface{}) interface{} {
		return Tuple{a,b}
	})
}

//Parser[B]
func (p Parser) As(b interface{}) Parser {
	return p.Slice().Map(func(i interface{}) interface{} {
		return b
	})
}

//Parser[C] A,B -> C
//事实上这个And会很少用，存在一个逻辑问题: p1 And p2 , p1解析产生出一个结果，p2解析出一个结果，这两个结果需要做一步什么操作才能继续下去？
func (p Parser) And(p2 LazyParser ,f mapFunc2) Parser {
	return p.Map2(p2,f)
}

//Parser[B >: A]
func (p Parser) Or(p2 LazyParser) Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if res.isSuccess() {
			return res
		} else {
			return p2()(state)
		}
	}
}

//Parser[B]
//省略左边
func (p Parser) SkipL(p2 LazyParser) Parser {
	return p.Map2(p2, func(a interface{}, b interface{}) interface{} {
		return b
	})
}

//Parser[A]
//省略右边
func (p Parser) SkipR(p2 LazyParser) Parser {
	return p.Map2(p2, func(a interface{}, b interface{}) interface{} {
		return a
	})
}

//Parser[List[A]]
//后面跟随的，0个或多个
func (p Parser) Sep(p2 LazyParser) Parser {
	var list = make([]interface{},0)
	return p.Sep1(p2).Or(MakeLazyParser(Succeed(list)))
}

//Parser[List[A]]
//后面跟随的，1个或多个
func (p Parser) Sep1(p2 LazyParser) Parser {
	loop := p2().SkipL(MakeLazyParser(p))
	many := MakeLazyParser(
		Many(MakeLazyParser(loop)))
	return p.Map2(many, func(a interface{}, b interface{}) interface{} {
		if list,ok:= b.([]interface{});ok{
			newList := []interface{}{a}
			newList = append(newList, list...)
			return newList
		} else {
			panic("parse value wrong type")
		}
	})
}