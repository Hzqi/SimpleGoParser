package parser

//实际上是一个 ParseState => ParseResult 的函数
type Parser func(state ParseState) ParseResult

//添加label，返回一个新的Parser(函数)
func (p Parser) Label(msg string) Parser {
	return func(state ParseState) ParseResult {
		return mapError(p(state), func(parseError ParseError) ParseError {
			return parseError.label(msg)
		})
	}
}

//添加scope
func (p Parser) Scope(msg string) Parser {
	return func(state ParseState) ParseResult {
		return mapError(p(state), func(parseError ParseError) ParseError {
			return *parseError.push(state.loc, msg)
		})
	}
}

func (p Parser) Attempt() Parser {
	return func(state ParseState) ParseResult {
		return uncommit(p(state))
	}
}

func (p Parser) Slice() Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if res.isSuccess() {
			return &Success{state.slice(res.getLength()), res.getLength()}
		} else {
			return res
		}
	}
}

func (p Parser) FlatMap(g func(interface{}) Parser) Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if res.isSuccess() {
			newRes := g(res.getSuccess())(state.advance(res.getLength()))
			afterCommit := addCommit(newRes, res.getLength() != 0)
			return advanceSuccess(afterCommit, res.getLength())
		} else {
			return res
		}
	}
}

func (p Parser) FlatMap2(p2 Parser, g func(interface{}, interface{}) Parser) Parser {
	return func(state ParseState) ParseResult {
		res1 := p(state)
		if res1.isSuccess() {
			state2 := state.advance(res1.getLength())
			res2 := p2(state2)
			if res2.isSuccess() {
				//两个都成功时
				newRes := g(res1.getSuccess(), res2.getSuccess())(state2.advance(res2.getLength()))
				afterCommit := addCommit(newRes, res2.getLength() != 0)
				return advanceSuccess(afterCommit, res2.getLength())
			} else {
				return res2
			}
		} else {
			return res1
		}
	}
}

//仿函数式，返回一个函数r， 函数r是执行完函数f后，用其结果再执行函数t
//类似于Haskell中的类型： andThen :: (a -> b) -> (b -> c) -> (a -> c)
func andThen(f func(interface{}) interface{}, t func(interface{}) Parser) func(interface{}) Parser {
	return func(i interface{}) Parser {
		return t(f(i))
	}
}

// (a -> b -> c) -> (c -> d) -> (a -> b -> d)
func andThen2(f func(interface{}, interface{}) interface{}, t func(interface{}) Parser) func(interface{}, interface{}) Parser {
	return func(i interface{}, i2 interface{}) Parser {
		return t(f(i, i2))
	}
}

func (p Parser) Map(f func(interface{}) interface{}) Parser {
	return p.FlatMap(andThen(f, Succeed))
}

func (p Parser) Map2(p2 Parser, f func(interface{}, interface{}) interface{}) Parser {
	return p.FlatMap2(p2, andThen2(f, Succeed))
}

func (p Parser) Product(p2 Parser) Parser {
	return p.FlatMap(func(i interface{}) Parser {
		return p2.Map(func(j interface{}) interface{} {
			return Tuple{i, j}
		})
	})
}

func (p Parser) As(b interface{}) Parser {
	return p.Slice().Map(func(i interface{}) interface{} {
		return b
	})
}

func (p Parser) And(p2 Parser) Parser {
	return p.FlatMap(func(i interface{}) Parser {
		return p2
	})
}

func (p Parser) Or(p2 Parser) Parser {
	return func(state ParseState) ParseResult {
		res := p(state)
		if !res.isSuccess() && !res.isSuccess() {
			return p2(state)
		} else {
			return res
		}
	}
}

func (p Parser) SkipL(p2 Parser) Parser {
	return p.Slice().Map2(p2, func(a interface{}, b interface{}) interface{} {
		return b
	})
}

func (p Parser) SkipR(p2 Parser) Parser {
	return p.Slice().Map2(p2, func(a interface{}, b interface{}) interface{} {
		return a
	})
}

func (p Parser) Sep(p2 Parser) Parser {
	return p.Sep1(p2).Or(Succeed(make([]interface{}, 5)))
}

func (p Parser) Sep1(p2 Parser) Parser {
	return p.Map2(Many(p2.SkipL(p)), func(a interface{}, b interface{}) interface{} {
		if list, ok := a.([]interface{}); ok {
			return append(list, b)
		} else {
			panic("parser result wrong type")
		}
	})
}

func (p Parser) ListOfN(n int) Parser {
	if n <= 0 {
		return Succeed(make([]interface{}, 5))
	} else {
		return p.Map2(p.ListOfN(n-1), func(a interface{}, b interface{}) interface{} {
			if list, ok := a.([]interface{}); ok {
				return append(list, b)
			} else {
				panic("parser result wrong type")
			}
		})
	}
}
