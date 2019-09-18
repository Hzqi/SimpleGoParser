package main

import (
	"fmt"
	"time"
)

type MyFun func(int)bool

func (f MyFun) label(name string) MyFun {
	return func(i int) bool {
		fmt.Printf("开始执行%s\n",name)
		return f(i)
	}
}

func (f MyFun) Or(f2 func()MyFun) MyFun {
	return func(i int) bool {
		res := f(i)
		if !res {
			fmt.Printf("不满足\n")
			return f2()(i)
		} else {
			return res
		}
	}
}

func a() MyFun {
	return func(i int) bool {
		return i == 1
	}
}

func b() MyFun {
	return func(i int) bool {
		return i == 2
	}
}

func c() MyFun {
	return func(i int) bool {
		return i == 3
	}
}

func main() {
	s := time.Now().UnixNano() / 1e6
	var aa = func() MyFun{return a().label("a") }
	var bb = func() MyFun{return b().label("b") }
	var cc = func() MyFun{return c().label("c") }

	f := aa().Or(bb).Or(cc)
	res := f(4)
	e := time.Now().UnixNano() / 1e6
	fmt.Printf("%s\n%dms    %d %d",res, e-s ,e,s)
}