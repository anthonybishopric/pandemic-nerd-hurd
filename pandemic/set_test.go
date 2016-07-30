package pandemic

import (
	"fmt"
	"testing"
)

func TestSetContains(t *testing.T) {
	s := Init([]string{"a", "b"})

	if c := s.Contains("a"); c != true {
		t.Fatalf("Should have contained %v got %v", "a", c)
	}

	if c := s.Contains("c"); c != false {
		t.Fatalf("Should not have contained %v got %v", "c", c)
	}

	s.Add("c")

	if c := s.Contains("c"); c != true {
		t.Fatalf("Should have contained %v got %v", "c", c)
	}

	s2 := Init([]string{"a"})

	s3 := Intersection(s, s2)

	if c := s3.Contains("a"); c != true {
		t.Fatalf("Should have contained %v got %v", "a", c)
	}

	fmt.Println(s3)
}
