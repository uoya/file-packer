package main

type Set map[string]bool

func NewSet() Set {
	return make(Set)
}

func (s Set) Add(value string) {
	s[value] = true
}

func (s Set) Values() []string {
	var values []string
	for k := range s {
		values = append(values, k)
	}
	return values
}
