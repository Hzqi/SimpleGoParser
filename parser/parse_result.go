package parser

type ParseResult interface {
	isSuccess() bool
	//success
	getSuccess() interface{}
	getLength() int
	//failure
	getError() ParseError
	isCommitted() bool
}

func extract(r ParseResult) Either {
	if r.isSuccess() {
		return &Right{r.getSuccess()}
	} else {
		return &Left{r.getError()}
	}
}

func uncommit(r ParseResult) ParseResult {
	if r.isSuccess() {
		return r
	} else {
		return &Failure{r.getError(), false}
	}
}

func addCommit(r ParseResult, committed bool) ParseResult {
	if r.isSuccess() {
		return r
	} else {
		return &Failure{r.getError(), r.isCommitted() || committed}
	}
}

func mapError(r ParseResult, fun func(ParseError) ParseError) ParseResult {
	if r.isSuccess() {
		return r
	} else {
		return &Failure{fun(r.getError()), r.isCommitted()}
	}
}

func advanceSuccess(r ParseResult, n int) ParseResult {
	if r.isSuccess() {
		return &Success{r.getSuccess(), r.getLength() + n}
	} else {
		return r
	}
}

type Success struct {
	value  interface{}
	length int
}

func (*Success) isSuccess() bool {
	return true
}

func (s *Success) getSuccess() interface{} {
	return s.value
}

func (s *Success) getLength() int {
	return s.length
}

func (*Success) getError() ParseError {
	return ParseError{}
}

func (*Success) isCommitted() bool {
	return false
}

type Failure struct {
	parseError ParseError
	committed  bool
}

func (*Failure) isSuccess() bool {
	return false
}

func (*Failure) getSuccess() interface{} {
	return nil
}

func (*Failure) getLength() int {
	return 0
}

func (f *Failure) getError() ParseError {
	return f.parseError
}

func (f *Failure) isCommitted() bool {
	return f.committed
}
