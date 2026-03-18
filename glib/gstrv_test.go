package glib

import (
	"testing"
)

func TestStrvBasic(t *testing.T) {
	s := NewStrv([]string{"hello", "world", "foo"})

	if s.Len() != 3 {
		t.Fatalf("expected len 3, got %d", s.Len())
	}

	v0, err := s.Index(0)
	if err != nil {
		t.Fatal(err)
	}
	if v0 != "hello" {
		t.Fatalf("expected hello, got %v", v0)
	}

	v1, err := s.Index(1)
	if err != nil {
		t.Fatal(err)
	}
	if v1 != "world" {
		t.Fatalf("expected world, got %v", v1)
	}

	v2, err := s.Index(2)
	if err != nil {
		t.Fatal(err)
	}
	if v2 != "foo" {
		t.Fatalf("expected foo, got %v", v2)
	}
}

func TestStrvMarshalRoundtrip(t *testing.T) {
	s := NewStrv([]string{"a", "b", "c"})

	gv, err := GValue(s)
	if err != nil {
		t.Fatal(err)
	}

	iface, err := gv.GoValue()
	if err != nil {
		t.Fatal(err)
	}

	s2, ok := iface.(*Strv)
	if !ok {
		t.Fatal("could not cast to *Strv")
	}

	if s2.Len() != 3 {
		t.Fatalf("expected len 3, got %d", s2.Len())
	}

	v, err := s2.Index(1)
	if err != nil {
		t.Fatal(err)
	}
	if v != "b" {
		t.Fatalf("expected b, got %v", v)
	}
}

func TestStrvEmpty(t *testing.T) {
	s := NewStrv([]string{})

	if s.Len() != 0 {
		t.Fatalf("expected len 0, got %d", s.Len())
	}

	strs := s.Strings()
	if len(strs) != 0 {
		t.Fatalf("expected empty slice, got %v", strs)
	}
}

func TestStrvOutOfBounds(t *testing.T) {
	s := NewStrv([]string{"only"})

	_, err := s.Index(-1)
	if err == nil {
		t.Fatal("expected error for negative index")
	}

	_, err = s.Index(1)
	if err == nil {
		t.Fatal("expected error for out of bounds index")
	}
}

func TestStrvStrings(t *testing.T) {
	input := []string{"alpha", "beta", "gamma"}
	s := NewStrv(input)

	result := s.Strings()
	if len(result) != len(input) {
		t.Fatalf("expected len %d, got %d", len(input), len(result))
	}
	for i, v := range input {
		if result[i] != v {
			t.Fatalf("index %d: expected %s, got %s", i, v, result[i])
		}
	}
}
