package parser

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLexImpl_Next(t *testing.T) {
	buf := bytes.NewBufferString("select A as \"My column\", B, C where A >= B")
	lex := NewLex(buf, []string{"select", "where", "as"}, []string{">", "=", ",", "*", "\"", "'"})
	for {
		token := lex.Next()
		if token == nil {
			break
		}
		println(fmt.Sprintf("[%v]=`%s`", token.Kind(), token.Text()))
	}
}
