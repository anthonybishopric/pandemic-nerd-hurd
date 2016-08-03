package pandemic

import (
	"testing"
)

func TestStack(t *testing.T) {
	s := NewStack()
	s.Push(1).Push(2)

	// Assert Peek LIFO
	if head := s.Peek(); head != 2 {
		t.Fatalf("A Peek should return the last element")
	}

	// Asset LIFO
	if head, _ := s.Pop(); head != 2 {
		t.Fatalf("Should have gotten last value pushed, got %v", head)
	}

	// Assert getting a single item
	if head, _ := s.Pop(); head != 1 {
		t.Fatalf("Should have gotten last value pushed, got %v", head)
	}

	// Assert we can handle empty stack
	if _, err := s.Pop(); err == nil {
		t.Fatalf("Should have gotten a 'stack empty' error, got %v", err)
	}

	// Assert we can handle empty stack again!
	if _, err := s.Pop(); err == nil {
		t.Fatalf("Should have gotten a 'stack empty' error, got %v", err)
	}
}
