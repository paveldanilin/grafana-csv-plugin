package parser

import (
	"io"
	"strconv"
	"text/scanner"
)

type TokenKind int

const (
	TokenKindKeyword = iota
	TokenKindOperator
	TokenKindIdentifier
	TokenKindNumeric
)

type Token struct {
	pos int
	line int
	text string
	kind TokenKind
}

func NewToken(kind TokenKind, text string, pos, line int) *Token {
	return &Token{
		pos:  pos,
		line: line,
		text: text,
		kind: kind,
	}
}

func (t *Token) Pos() int {
	return t.pos
}

func (t *Token) Line() int {
	return t.line
}

func (t *Token) Text() string {
	return t.text
}

func (t *Token) Kind() TokenKind {
	return t.kind
}

type Lex interface {
	Next() *Token
}

type LexImpl struct {
	scanner    scanner.Scanner
	keywords  []string
	operators []string
}

func NewLex(reader io.Reader, keywords []string, operators []string) Lex {
	lex := &LexImpl{
		keywords:  keywords,
		operators: operators,
	}
	lex.scanner.Init(reader)
	return lex
}

func (lex *LexImpl) Next() *Token {
	tok := lex.scanner.Scan()

	if tok == scanner.EOF {
		return nil
	}

	pos := lex.scanner.Pos()
	text := lex.scanner.TokenText()

	return NewToken(lex.resolveKind(text), text, pos.Column, pos.Line)
}

func (lex *LexImpl) resolveKind(text string) TokenKind {
	if lex.isKeyword(text) {
		return TokenKindKeyword
	}
	if lex.isOperator(text) {
		return TokenKindOperator
	}
	if lex.isNumeric(text) {
		return TokenKindNumeric
	}
	return TokenKindIdentifier
}

func (lex *LexImpl) isKeyword(text string) bool {
	for _, kw := range lex.keywords {
		if kw == text {
			return true
		}
	}
	return false
}

func (lex *LexImpl) isOperator(text string) bool {
	for _, op := range lex.operators {
		if op == text {
			return true
		}
	}
	return false
}

func (lex *LexImpl) isNumeric(text string) bool {
	_, err := strconv.ParseFloat(text, 64)
	return err == nil
}
