package ldapfilter

import (
	"reflect"
	"testing"
)

func TestFilterAnd(t *testing.T) {
	input := Input{
		"foo": {"bar", "baz"},
	}
	filter := AndFilter{}
	filter.Append(&AndFilter{})
	res := filter.Match(input)
	assert(t, true, res)
}

func TestFilterOr(t *testing.T) {
	input := Input{
		"foo": {"bar", "baz"},
	}
	filter := &OrFilter{
		Children: []Filter{
			&AndFilter{},
			&AndFilter{},
		},
	}
	res := filter.Match(input)
	assert(t, true, res)
}

func TestFilterEquality(t *testing.T) {
	filter := EqualityFilter{
		Key:   "foo",
		Value: "baz",
	}

	t.Run("should match", func(t *testing.T) {
		input := Input{
			"foo": {"bar", "baz"},
		}
		res := filter.Match(input)
		assert(t, true, res)
	})

	t.Run("should not match", func(t *testing.T) {
		input := Input{
			"foo": {"bar"},
		}
		res := filter.Match(input)
		assert(t, false, res)
	})
}

func TestFilter(t *testing.T) {

	testData := []struct {
		expr string
		ok   bool
		inp  Input
	}{
		{
			expr: "(name=Jon)",
			inp:  Input{"name": {"Jon"}},
			ok:   true,
		},
		{
			expr: "|(name=Jon)(name=Foo)",
			inp:  Input{"name": {"Jon"}},
			ok:   true,
		},
		{
			expr: "|(name=Jon)(alt=Foo)",
			inp:  Input{"alt": {"Foo"}},
			ok:   true,
		},
	}

	for _, test := range testData {
		t.Run("", func(t *testing.T) {
			filter, err := NewParser(test.expr).Parse()
			assert(t, nil, err)
			assert(t, test.ok, filter.Match(test.inp))
		})
	}
}

// helper

func assert(t *testing.T, want interface{}, have interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, have) {
		t.Errorf("Assertion failed for %s\n\twant:\t%+v\n\thave:\t%+v", t.Name(), want, have)
	}
}
