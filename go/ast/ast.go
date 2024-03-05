package ast

import (
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	out := strings.Builder{}
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
