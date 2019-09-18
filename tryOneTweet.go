package main

import (
	"SimpleGoParser/parser"
	"fmt"
)

type OneTweet interface{}

type OneA struct {
	ones []OneTweet
}
type OneB struct {
	ones []OneTweet
}
type OneS struct {
	str string
}

func root() parser.Parser {
	return parserA().Or(parserB()).Or(parserS())
}
func parserA() parser.Parser {
	return parser.Surround(parser.Str("{"), parser.Str("}"), parser.Many(root())).
		Map(func(i interface{}) interface{} {
			return OneA{i.([]OneTweet)}
		})
}

func parserB() parser.Parser {
	return parser.Surround(parser.Str("("), parser.Str(")"), parser.Many(root())).
		Map(func(i interface{}) interface{} {
			return OneB{i.([]OneTweet)}
		})
}

func parserS() parser.Parser {
	return parser.Surround(parser.Str("<"), parser.Str(">"), parser.Many(parser.CharNotIn([]rune("(){}<>")))).
		Map(func(i interface{}) interface{} {
			return OneS{string(i.([]rune))}
		})
}
func main() {
	var str = "{()(){}{}<abcd><1234>}"
	state := parser.NewParseState(*parser.NewLocationDefault(str))
	res := root()(state)
	fmt.Printf("%v\n", res)
}
