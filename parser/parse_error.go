package parser

import (
	"fmt"
	"sort"
	"strings"
)

type ParseError struct {
	stack []LocStr
}

type LocStr struct {
	loc Location
	str string
}

//为实现Tuple排序的别名
type LocStrIntf []LocStr

func (tf LocStrIntf) Len() int {
	return len(tf)
}

func (tf LocStrIntf) Less(i, j int) bool {
	ti := tf[i]
	tj := tf[j]
	return ti.loc.offset < tj.loc.offset
}

func (tf LocStrIntf) Swap(i, j int) {
	tf[i], tf[j] = tf[j], tf[i]
}

func (pe *ParseError) push(location Location, msg string) *ParseError {
	stack := append(pe.stack, LocStr{location, msg})
	pe.stack = stack
	return pe
}

func (pe *ParseError) label(s string) ParseError {
	if loc, ok := pe.latestLoc(); ok {
		tuple := LocStr{loc, s}
		return ParseError{[]LocStr{tuple}}
	} else {
		return *pe
	}
}

func (pe *ParseError) latest() (LocStr, bool) {
	length := len(pe.stack)
	if length > 0 {
		return pe.stack[length-1], true
	} else {
		return LocStr{}, false
	}
}

func (pe *ParseError) latestLoc() (Location, bool) {
	if tuple, ok := pe.latest(); ok {
		return tuple.loc, true
	} else {
		return Location{}, false
	}
}

func (pe *ParseError) collapseStack(s []LocStr) []LocStr {
	group := groupBy(pe.stack)
	list := make([]LocStr, len(group))
	for k, v := range group {
		list = append(list, LocStr{k, v})
	}
	sort.Sort(LocStrIntf(list))
	return list
}

func (pe *ParseError) formatLoc(loc Location) string {
	return fmt.Sprintf("%d.%d", loc.line(), loc.col())
}

func (pe *ParseError) String() string {
	if len(pe.stack) == 0 {
		return "no error message"
	} else {
		collapsed := pe.collapseStack(pe.stack)
		length := len(collapsed)
		context := ""
		if length > 0 {
			context = fmt.Sprintf("\n\n%s \n%s", collapsed[length-1].loc.currentLine(), collapsed[length-1].loc.columnCaret())
		}
		list := make([]string, length+1)
		for i, v := range collapsed {
			str := fmt.Sprintf("%d.%d except:%s", v.loc.line(), v.loc.col(), v.str)
			list[i] = str
		}
		list[length] = context
		return strings.Join(list, "\n")
	}
}

//仿照Scala实现的groupBy，但是泛型就不行了
func groupBy(list []LocStr) map[Location]string {
	group := make(map[Location]string)
	for _, t := range list {
		loc := t.loc
		str := t.str
		if fromMap, ok := group[loc]; ok {
			fromMap += ";" + str
			group[loc] = fromMap
		} else {
			fromMap = str
			group[loc] = fromMap
		}
	}
	return group
}
