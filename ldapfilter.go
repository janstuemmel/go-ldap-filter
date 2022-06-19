package ldapfilter

type Input = map[string][]interface{}

type Filter interface {
	Match(input Input) bool
	Append(filter Filter)
}

// AND

type AndFilter struct {
	Type     string
	Children []Filter
}

func NewAndFilter() *AndFilter {
	return &AndFilter{Type: "and"}
}

func (f *AndFilter) Match(input Input) bool {
	for _, filter := range f.Children {
		if !filter.Match(input) {
			return false
		}
	}
	return true
}

func (f *AndFilter) Append(filter Filter) {
	f.Children = append(f.Children, filter)
}

// OR

type OrFilter struct {
	Type     string
	Children []Filter
}

func NewOrFilter() *OrFilter {
	return &OrFilter{Type: "or"}
}

func (f *OrFilter) Match(input Input) bool {
	if len(f.Children) == 0 {
		return true
	}
	for _, filter := range f.Children {
		if filter.Match(input) {
			return true
		}
	}
	return false
}

func (f *OrFilter) Append(filter Filter) {
	f.Children = append(f.Children, filter)
}

// EQUALITY

type EqualityFilter struct {
	Type  string
	Key   string
	Value string
}

func NewEqualityFilter() *EqualityFilter {
	return &EqualityFilter{Type: "equality"}
}

func (f *EqualityFilter) Match(input Input) bool {
	if values, ok := input[f.Key]; ok {
		for _, value := range values {
			if value == f.Value {
				return true
			}
		}
	}
	return false
}

func (f *EqualityFilter) Append(filter Filter) {}
