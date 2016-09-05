package pandemic

import (
	"testing"
)

type testStringable string

func (t testStringable) String() string {
	return string(t)
}

func TestSetContains(t *testing.T) {
	s := Init(testStringable("a"), testStringable("b"))

	if c := s.Contains(testStringable("a")); c != true {
		t.Fatalf("Should have contained %v got %v", "a", c)
	}

	if c := s.Contains(testStringable("c")); c != false {
		t.Fatalf("Should not have contained %v got %v", "c", c)
	}

	s.Add(testStringable("c"))

	if c := s.Contains(testStringable("c")); c != true {
		t.Fatalf("Should have contained %v got %v", "c", c)
	}

	s2 := Init(testStringable("a"))

	s3 := Intersection(s, s2)

	if c := s3.Contains(testStringable("a")); c != true {
		t.Fatalf("Should have contained %v got %v", "a", c)
	}
}

func TestSet_Members(t *testing.T) {
	s := Set{"foo": struct{}{}, "bar": struct{}{}}
	if members := s.Members(); members[0] != "bar" || members[1] != "foo" || len(members) != 2 {
		t.Fatalf("Returns members of a set in sorted order %v", members)
	}
}
