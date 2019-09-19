package parser

//位置信息， 一个输入源（字符串），一个下标开始位置
type Location struct {
	input  string //字符串是不会变的
	offset int
}

func NewLocation(input string, offset int) Location {
	return Location{input,offset}
}