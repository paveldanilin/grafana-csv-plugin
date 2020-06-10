package parser

import (
	"bytes"
)

type Statement interface {
	Supports(first *Token) bool
}

type Node struct {
	token *Token
	left *Node
	right *Node
}

type Rule interface {
	Supports(first *Token) bool
}

type RuleSelect struct {

}

func (r *RuleSelect) Supports(first *Token) bool {
	return first.Kind() == TokenKindKeyword && first.Text() == "select"
}

type Parser interface {
	Parse(program string) *Node
}

type SqlParser struct {
	rules []Rule
}

func NewSqlParser() Parser {
	rules := make([]Rule, 0)
	rules = append(rules, &RuleSelect{})
	return &SqlParser{
		rules: rules,
	}
}

func (p *SqlParser) Parse(program string) *Node {
	lex := NewLex(bytes.NewBufferString(program),  []string{"select", "where", "as"}, []string{">", "=", ",", "*", "\"", "'"})

	rule := p.lookup(lex.Next())
	if rule == nil {
		// Not found rule
		return nil
	}

	return nil
}

func (p *SqlParser) lookup(first *Token) Rule {
	for _, r := range p.rules {
		if r.Supports(first) {
			return r
		}
	}
	return nil
}
