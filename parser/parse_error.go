package parser

import (
	"fmt"
	"strings"
)

type ParseError struct {
	currentLine string
	msg string
 	line int
	col int
}

func NewParseError(currentLine , msg string, line,col int) ParseError {
	return ParseError{currentLine,msg,line,col}
}

func getLineAndCol(input string, offset int) (string,int,int) {
	currInput := input[offset : ]
	currInputs := strings.Split(currInput,"\n")
	lineCount := len(currInputs)
	line := currInputs[lineCount - 1]
	col := strings.LastIndex(currInput,"\n")
	if col == -1 {
		col = offset + 1
	} else {
		col = offset - col
	}
	return line,lineCount,col
}

func (pe *ParseError) Error() string {
	linemsg := fmt.Sprintf(" line:%d col:%d %s",pe.line,pe.col, pe.msg)
	line := fmt.Sprintf("\t%s",pe.currentLine)
	blank := strings.Repeat("",pe.col)
	mark := fmt.Sprintf("\t%s^",blank)
	return strings.Join([]string{linemsg,line,mark},"\n")
}