package test

import (
	"SimpleGoParser/example"
	"SimpleGoParser/parser"
	"testing"
)

/*
{{}{}{}} -> T {
				[]T{ T{[]{}}, T{[]{}}, T{[]{}}
			}
 */
type T struct {
	children []T
}

func root() parser.Parser {
	return parser.Surround(
		parser.MakeLazyParser(parser.Str("{")),
		parser.MakeLazyParser(parser.Str("}")),
		manyRoot).
		Map(func(i interface{}) interface{} {
			if list,ok := i.([]interface{}); ok {
				var ts []T
				for _,v := range list {
					switch v.(type) {
						case T : ts = append(ts, v.(T))
						case []T : ts = append(ts, T{v.([]T)})
					}
				}
				return T{ts}
			} else {
				panic("wrong type")
			}
	})
}

func manyRoot() parser.Parser {
	return parser.Many(root)
}

func TestExample(t *testing.T) {
	input := "{{}{}{}}"
	t.Log(RunParser(root(),input))
}

func TestJson(t *testing.T) {
	input := `
	 {
		"name":"jacky",
		"age":123,
		"gender":true,
		"nulls": null,
 		"object":{
            "aaa":"aaa",
            "bbb":"bbb"
		},
		"arrays":[1,   2,   3,     4     ],
		"nulls": null,
		"str":"abc\"aaa\"",
		"cn_Name":"中文名"
	 }
`
	res,err := parser.RunParser(example.RootWithBlank(), input)
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(res)
	}
}