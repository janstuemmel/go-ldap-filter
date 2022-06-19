package ldapfilter

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	lexer := NewLexer(strings.NewReader("&hello() 1=|"))
	assert(t, lexer.Lex().Type, AND)
	assert(t, lexer.Lex().Type, IDENT)
	assert(t, lexer.Lex().Type, OPEN)
	assert(t, lexer.Lex().Type, CLOSE)
	assert(t, lexer.Lex().Type, ILLEGAL)
	assert(t, lexer.Lex().Type, ILLEGAL)
	assert(t, lexer.Lex().Type, EQUAL)
	assert(t, lexer.Lex().Type, OR)
	assert(t, lexer.Lex().Type, EOF)
	assert(t, lexer.Lex().Type, EOF)
}

func TestParseError(t *testing.T) {
	exprs := []string{
		"",
		"hello",
		"(",
		")",
		"foo=bar",
		"&()",
		"&(foo=)",
		"&(=bar)",
		// "&)", // TODO
	}
	for _, expr := range exprs {
		t.Run(fmt.Sprintf("parse '%s'", expr), func(t *testing.T) {
			_, err := NewParser(expr).Parse()
			assert(t, ErrParser, err)
		})
	}
}

func TestParse(t *testing.T) {

	t.Run("parse simple expression", func(t *testing.T) {
		f, err := NewParser("(foo=bar)").Parse()
		assert(t, nil, err)
		assert(t, &EqualityFilter{
			Type:  "equality",
			Key:   "foo",
			Value: "bar",
		}, f)
	})

	t.Run("parse filter", func(t *testing.T) {
		f, err := NewParser("&(foo=bar)(bar=baz)").Parse()
		assert(t, nil, err)
		assert(t, &AndFilter{
			Type: "and",
			Children: []Filter{
				&EqualityFilter{
					Type:  "equality",
					Key:   "foo",
					Value: "bar",
				},
				&EqualityFilter{
					Type:  "equality",
					Key:   "bar",
					Value: "baz",
				},
			},
		}, f)
	})

	t.Run("parse complex filter", func(t *testing.T) {
		f, err := NewParser("|(&(foo=bar)(bar=baz))(&(foo=bar)(bar=baz))").Parse()
		assert(t, nil, err)
		assert(t, &OrFilter{
			Type: "or",
			Children: []Filter{
				&AndFilter{
					Type: "and",
					Children: []Filter{
						&EqualityFilter{
							Type:  "equality",
							Key:   "foo",
							Value: "bar",
						},
						&EqualityFilter{
							Type:  "equality",
							Key:   "bar",
							Value: "baz",
						},
					},
				},
				&AndFilter{
					Type: "and",
					Children: []Filter{
						&EqualityFilter{
							Type:  "equality",
							Key:   "foo",
							Value: "bar",
						},
						&EqualityFilter{
							Type:  "equality",
							Key:   "bar",
							Value: "baz",
						},
					},
				},
			},
		}, f)
	})
}
