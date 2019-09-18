package parser

//仿照Scala的List.groupBy来构建一个map
//没有泛型，配合interface{}的话感觉写不出来

// Either , Left, Right
//=====================================
type Either interface {
	IsLeft() bool
	IsRight() bool

	GetRight() interface{}
	GetLeft() interface{}
}
type Right struct {
	Value interface{}
}

func (*Right) IsLeft() bool {
	return false
}

func (*Right) IsRight() bool {
	return true
}

func (r *Right) GetRight() interface{} {
	return r.Value
}

func (r *Right) GetLeft() interface{} {
	return nil
}

type Left struct {
	Value interface{}
}

func (*Left) IsLeft() bool {
	return true
}

func (*Left) IsRight() bool {
	return false
}

func (*Left) GetRight() interface{} {
	return nil
}

func (l *Left) GetLeft() interface{} {
	return l.Value
}

//=====================================

//Tuple
//=====================================
type Tuple struct {
	One interface{}
	Two interface{}
}
