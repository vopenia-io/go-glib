package glib

import (
	"testing"
)

func TestArrayBasic(t *testing.T) {
	arr, err := NewArray([]interface{}{42, "hello", 3.14})
	if err != nil {
		t.Fatal(err)
	}

	if arr.Len() != 3 {
		t.Fatalf("expected len 3, got %d", arr.Len())
	}

	v0, err := arr.Index(0)
	if err != nil {
		t.Fatal(err)
	}
	if v0.(int) != 42 {
		t.Fatalf("expected 42, got %v", v0)
	}

	v1, err := arr.Index(1)
	if err != nil {
		t.Fatal(err)
	}
	if v1.(string) != "hello" {
		t.Fatalf("expected hello, got %v", v1)
	}

	v2, err := arr.Index(2)
	if err != nil {
		t.Fatal(err)
	}
	if v2.(float64) != 3.14 {
		t.Fatalf("expected 3.14, got %v", v2)
	}
}

func TestArrayMarshalRoundtrip(t *testing.T) {
	arr, err := NewArray([]interface{}{1, "two", 3.0})
	if err != nil {
		t.Fatal(err)
	}

	gv, err := GValue(arr)
	if err != nil {
		t.Fatal(err)
	}

	iface, err := gv.GoValue()
	if err != nil {
		t.Fatal(err)
	}

	arr2, ok := iface.(*Array)
	if !ok {
		t.Fatal("could not cast to *Array")
	}

	if arr2.Len() != 3 {
		t.Fatalf("expected len 3, got %d", arr2.Len())
	}

	v, err := arr2.Index(1)
	if err != nil {
		t.Fatal(err)
	}
	if v.(string) != "two" {
		t.Fatalf("expected two, got %v", v)
	}
}

func TestArrayNested(t *testing.T) {
	row0, err := NewArray([]interface{}{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}

	row1, err := NewArray([]interface{}{4, 5, 6})
	if err != nil {
		t.Fatal(err)
	}

	matrix, err := NewArray([]interface{}{row0, row1})
	if err != nil {
		t.Fatal(err)
	}

	if matrix.Len() != 2 {
		t.Fatalf("expected len 2, got %d", matrix.Len())
	}

	iface, err := matrix.Index(0)
	if err != nil {
		t.Fatal(err)
	}

	inner, ok := iface.(*Array)
	if !ok {
		t.Fatal("could not cast to *Array")
	}

	v, err := inner.Index(2)
	if err != nil {
		t.Fatal(err)
	}
	if v.(int) != 3 {
		t.Fatalf("expected 3, got %v", v)
	}
}

func TestArrayOutOfBounds(t *testing.T) {
	arr, err := NewArray([]interface{}{1})
	if err != nil {
		t.Fatal(err)
	}

	_, err = arr.Index(-1)
	if err == nil {
		t.Fatal("expected error for negative index")
	}

	_, err = arr.Index(1)
	if err == nil {
		t.Fatal("expected error for out of bounds index")
	}
}

func TestArrayEmpty(t *testing.T) {
	arr, err := NewArray([]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	if arr.Len() != 0 {
		t.Fatalf("expected len 0, got %d", arr.Len())
	}

	vals, err := arr.Values()
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 0 {
		t.Fatalf("expected empty slice, got %v", vals)
	}
}
