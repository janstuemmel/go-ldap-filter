/*
Package ldapfilter provides functions for parsing
ldap filters and validating ldap entries.
*/
package ldapfilter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
)

// ErrParser error when parsing error happend
var ErrParser = errors.New("parser Error")

type (

	// Token emited by the lexer
	Token struct {
		Type int
		Lit  string
	}

	// Lexer used for lexing ldap filter strings
	Lexer struct {
		r *bufio.Reader
	}

	// Parser used for parsing ldap filter strings
	Parser struct {
		lexer  *Lexer
		pos    int     // keep track of token position
		tokens []Token // keep track of tokens
	}
)

const (
	// ILLEGAL Unknown token occured
	ILLEGAL int = iota
	// EOF End of file
	EOF
	// IDENT is the left or right of a ldap criteria
	IDENT

	// OPEN bracket
	OPEN
	// CLOSE bracket
	CLOSE

	// AND filter
	AND
	// OR filter
	OR
	// NEG negation filter
	NEG
	// EQUAL equality filter
	EQUAL
)

// NewLexer initializes a lexer
func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(reader)}
}

func (l *Lexer) lexIdent() Token {
	var buf bytes.Buffer

	for {
		r, _, err := l.r.ReadRune()

		if err != nil {
			break
		}

		if unicode.IsLetter(r) {
			buf.WriteRune(r)
		} else {
			l.r.UnreadRune()
			break
		}
	}

	return Token{IDENT, buf.String()}
}

// Lex emits next token
func (l *Lexer) Lex() Token {
	r, _, err := l.r.ReadRune()

	if err != nil {
		return Token{EOF, ""}
	}

	if unicode.IsLetter(r) {
		l.r.UnreadRune()
		return l.lexIdent()
	}

	switch r {
	default:
		return Token{ILLEGAL, ""}
	case '(':
		return Token{OPEN, string(r)}
	case ')':
		return Token{CLOSE, string(r)}
	case '&':
		return Token{AND, string(r)}
	case '|':
		return Token{OR, string(r)}
	case '!':
		return Token{NEG, string(r)}
	case '=':
		return Token{EQUAL, string(r)}
	}
}

// NewParser initializes a new Parser
func NewParser(str string) *Parser {
	return &Parser{
		lexer: NewLexer(strings.NewReader(str)),
	}
}

func (p *Parser) scan() (token Token) {
	// return buffered token from unscan
	if p.pos < len(p.tokens) {
		token = p.tokens[p.pos]
		p.pos++
		return
	}

	// return next lexer token
	token = p.lexer.Lex()
	p.tokens = append(p.tokens, token)
	p.pos++
	return
}

func (p *Parser) unscan(i int) {
	if p.pos-i <= 0 {
		p.pos = 0
		return
	}
	p.pos = p.pos - i
}

func (p *Parser) parseExpression() (Filter, error) {
	filter := NewEqualityFilter()

	tok := p.scan()
	if tok.Type != OPEN {
		return filter, ErrParser
	}

	if tok = p.scan(); tok.Type == IDENT {
		filter.Key = tok.Lit
	} else {
		return filter, ErrParser
	}

	if tok = p.scan(); tok.Type != EQUAL {
		return filter, ErrParser
	}

	if tok = p.scan(); tok.Type == IDENT {
		filter.Value = tok.Lit
	} else {
		return filter, ErrParser
	}

	if tok = p.scan(); tok.Type != CLOSE {
		return filter, ErrParser
	}

	return filter, nil
}

func (p *Parser) parseFilter() (Filter, error) {
	var filter Filter
	tok := p.scan()

	if tok.Type == AND || tok.Type == OR {

		switch tok.Type {
		default:
			filter = NewAndFilter()
		case OR:
			filter = NewOrFilter()
		}

		for {
			tok := p.scan()
			if tok.Type == EOF || tok.Type == CLOSE {
				return filter, nil
			}

			f, err := p.parseFilter()
			if err != nil {
				return filter, ErrParser
			}

			filter.Append(f)
		}

	} else {
		p.unscan(2)

		f, err := p.parseExpression()
		if err != nil {
			return filter, ErrParser
		}

		filter = f
	}

	return filter, nil
}

// Parse parses the ldap filter string
func (p *Parser) Parse() (Filter, error) {
	return p.parseFilter()
}
