package parser

// ParseResult[A]
type ParseResult interface {
	isSuccess() bool
	extract() (interface{}, error) //(A, error) 成功有返回值，失败就nil，有err
	//success
	getSuccess() interface{} //A
	getLength() int
	//failure
	getError() ParseError
}

//Success[A]
type Success struct {
	value  interface{} //A
	length int         //读取input的长度
}

func (*Success) isSuccess() bool { return true }

func (s *Success) extract() (interface{}, error) {
	return s.value, nil
}

func (s *Success) getSuccess() interface{} { return s.value }

func (s *Success) getLength() int { return s.length }

func (*Success) getError() ParseError { return ParseError{} }

//Failure[A]
type Failure struct {
	parseError ParseError
}

func (*Failure) isSuccess() bool { return false }

func (f *Failure) extract() (interface{}, error) {
	return nil, &f.parseError
}

func (*Failure) getSuccess() interface{} { return nil }

func (*Failure) getLength() int { return 0 }

func (f *Failure) getError() ParseError {
	return f.parseError
}