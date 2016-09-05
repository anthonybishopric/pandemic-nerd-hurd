package pandemic

import (
	"sort"
	"strings"
)

type Stringable interface {
	String() string
}

type Set map[string]struct{}

func Init(ks ...Stringable) Set {
	s := Set{}
	for _, k := range ks {
		s[k.String()] = struct{}{}
	}
	return s
}

func (s Set) Contains(k Stringable) bool {
	_, ok := s[k.String()]
	return ok
}

func (s Set) Add(k Stringable) Set {
	s[k.String()] = struct{}{}
	return s
}

func (s Set) Remove(k Stringable) (Set, bool) {
	if _, ok := s[k.String()]; !ok {
		return s, false
	}
	delete(s, k.String())
	return s, true
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Members() []string {
	ret := []string{}
	for k, _ := range s {
		ret = append(ret, k)
	}
	sort.Sort(SortedNames(ret))
	return ret
}

type SortedNames []string

func (s SortedNames) Less(i, j int) bool {
	return strings.Compare(s[i], s[j]) < 0
}

func (s SortedNames) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedNames) Len() int {
	return len(s)
}

func Intersection(s1 Set, s2 Set) Set {
	s3 := Set{}
	for k, _ := range s1 {
		if s2.Contains(stringer(k)) {
			s3.Add(stringer(k))
		}
	}
	return s3
}

type stringer string

func (s stringer) String() string {
	return string(s)
}
