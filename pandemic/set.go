package pandemic

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Stringable interface {
	String() string
}

type Set map[Stringable]struct{}

func Init(ks ...Stringable) Set {
	s := Set{}
	for _, k := range ks {
		s[k] = struct{}{}
	}
	return s
}

func (s Set) Contains(k Stringable) bool {
	_, ok := s[k]
	return ok
}

func (s Set) Add(k Stringable) Set {
	s[k] = struct{}{}
	return s
}

func (s Set) Remove(k Stringable) (Set, bool) {
	if _, ok := s[k]; !ok {
		return s, false
	}
	delete(s, k)
	return s, true
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Members() []Stringable {
	ret := []Stringable{}
	for k, _ := range s {
		ret = append(ret, k)
	}
	sort.Sort(SortedNames(ret))
	return ret
}

type SortedNames []Stringable

func (s SortedNames) Less(i, j int) bool {
	return strings.Compare(s[i].String(), s[j].String()) < 0
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
		if s2.Contains(k) {
			s3.Add(k)
		}
	}
	return s3
}

func (set Set) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]struct{}{}
	for k, _ := range set {
		toMarshal[k.String()] = struct{}{}
	}
	return json.Marshal(toMarshal)
}

type stringer string

func (s stringer) String() string {
	return string(s)
}

func (set Set) UnmarshalJSON(data []byte) error {
	var s map[string]struct{}
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("set should be a map of string to struct{}, got %s", data)
	}
	for k, _ := range s {
		set[stringer(k)] = struct{}{}
	}
	return nil
}
