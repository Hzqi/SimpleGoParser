package example

import (
	"SimpleGoParser/parser"
)

type JValue interface {}

type JObj struct {
	JValue //Use for implement
	Vals map[string]JValue
}
type JAry struct {
	JValue //Use for implement
	Vals []JValue
}
type JString struct {
	JValue //Use for implement
	Val string
}
type JNumber struct {
	JValue //Use for implement
	Val float64
}
type JBool struct {
	JValue //Use for implement
	Val bool
}
type JNull struct {
	JValue //Use for implement
}

//读null
func jNull() parser.Parser {
	return parser.Str("null").As(JNull{})
}

//读浮点数，也就是number
func jDouble() parser.Parser{
	return parser.Double().Map(func(i interface{}) interface{} {
		double := i.(float64)
		return JNumber{Val:double}
	})
}

//读字符串字面量
func jLitString() parser.Parser {
	return parser.EscapedQuoted().Map(func(i interface{}) interface{} {
		str := i.(string)
		return JString{Val:str}
	})
}

//读true
func jTrue() parser.Parser {
	return parser.Str("true").As(JBool{Val:true})
}

//读false
func jFalse() parser.Parser {
	return parser.Str("false").As(JBool{Val:false})
}

//读字面量
func lit() parser.Parser {
	return jNull().
		Or(jDouble).
		Or(jLitString).
		Or(jTrue).
		Or(jFalse).
		Scope("literal")
}

//读值，包括字面量、对象、数组
func jvalue() parser.Parser {
	return lit().Or(jObj).Or(jAry)
}

//读:成一个键值对
func keyval() parser.Parser {
	afterColon := parser.Token(parser.LazyStr(":")).SkipL(jvalue)
	return parser.Quoted().Product(
		parser.MakeLazyParser(afterColon))
}

//读数组内单元(后面跟，)
func arrayUnit() parser.Parser {
	return parser.Token(jvalue).Sep(parser.LazyStr(","))
}

//读数组
func jAry() parser.Parser {
	return parser.Surround(parser.LazyStr("["),parser.LazyStr("]"), arrayUnit).Map(func(i interface{}) interface{} {
			list := i.([]interface{})
			array := make([]JValue,len(list))
			for i,v := range list {
				array[i] = v
			}
			return array
		}).Scope("array")
}

//读对象内单元(后面跟，)
func objectUnit() parser.Parser {
	return parser.Token(keyval).Sep(parser.LazyStr(","))
}

//读对象
func jObj() parser.Parser {
	return parser.Surround(parser.LazyStr("{"), parser.LazyStr("}"), objectUnit).Map(func(i interface{}) interface{} {
		pairs := i.([]interface{})
		kvmap := make(map[string]JValue)
		for _, v := range pairs {
			t := v.(parser.Tuple)
			key := t.One.(string)
			value := t.Two.(JValue)
			kvmap[key] = value
		}
		return JObj{Vals:kvmap}
	})
}

//根parser
func root() parser.Parser {
	return jObj().Or(jAry)
}

//前面带空格的根
func RootWithBlank() parser.Parser {
	return parser.Whitespace().SkipL(root)
}