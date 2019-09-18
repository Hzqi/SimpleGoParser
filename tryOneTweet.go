package main

import (
	"SimpleGoParser/parser"
	"fmt"
)
type OneA struct {
	ones []OneA
}
func root() parser.Parser {
	return parserA()
}
func parserA() parser.Parser {
	return parser.Surround(parser.Str("{"), parser.Str("}"), parser.Many(root)).
		Map(func(i interface{}) interface{} {
			return OneA{i.([]OneA)}
		})
}
func main() {
	var str = "{}"
	resf := root()
	fmt.Printf("%v\n",resf)
	res := parser.RunParser(root(),str)
	//stringer := res.GetLeft().(parser.ParseError)
	fmt.Printf("%v\n %s\n", res/*, (&stringer).String()*/)
}
