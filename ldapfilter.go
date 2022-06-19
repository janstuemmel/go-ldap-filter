package ldapfilter

type Input = map[string][]interface{}

type Filter interface {
	Match(input Input) bool
	Append(filter Filter)
}

// AndFilter filters entries by and operation
type AndFilter struct {
	Type     string
	Children []Filter
}

// NewAndFilter inititalizes a new AndFilter
func NewAndFilter() *AndFilter {
	return &AndFilter{Type: "and"}
}

// Match matches entry input
func (f *AndFilter) Match(input Input) bool {
	for _, filter := range f.Children {
		if !filter.Match(input) {
			return false
		}
	}
	return true
}

// Append appends a new filter
func (f *AndFilter) Append(filter Filter) {
	f.Children = append(f.Children, filter)
}

// OrFilter filters entries by or operation
type OrFilter struct {
	Type     string
	Children []Filter
}

// NewOrFilter inititalizes a new OrFilter
func NewOrFilter() *OrFilter {
	return &OrFilter{Type: "or"}
}

// Match matches entry input
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

// Append appends a new filter
func (f *OrFilter) Append(filter Filter) {
	f.Children = append(f.Children, filter)
}

// EqualityFilter filters entries by equality
type EqualityFilter struct {
	Type  string
	Key   string
	Value string
}

// NewEqualityFilter inititalizes a new EqualityFilter
func NewEqualityFilter() *EqualityFilter {
	return &EqualityFilter{Type: "equality"}
}

// Match matches entry input
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

// Append does not appand anything here,
// because EqualityFilter is already a leaf
func (f *EqualityFilter) Append(filter Filter) {}
