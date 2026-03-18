package glib

// #include "glib.go.h"
//
// static void _gvalue_clear_func(gpointer data) {
//     GValue *v = (GValue *)data;
//     if (G_IS_VALUE(v)) { g_value_unset(v); }
// }
//
// static GArray *_g_array_new_gvalues(guint reserved) {
//     GArray *arr = g_array_sized_new(FALSE, TRUE, sizeof(GValue), reserved);
//     g_array_set_clear_func(arr, _gvalue_clear_func);
//     return arr;
// }
//
// static void _g_array_append_gvalue(GArray *array, const GValue *value) {
//     GValue copy = G_VALUE_INIT;
//     g_value_init(&copy, G_VALUE_TYPE(value));
//     g_value_copy(value, &copy);
//     g_array_append_vals(array, &copy, 1);
// }
//
// static GValue *_g_array_index_gvalue(GArray *array, guint i) {
//     return &g_array_index(array, GValue, i);
// }
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

var TYPE_ARRAY Type = Type(C.G_TYPE_ARRAY) // is function g_array_get_type() inside macro, can't be const in Go

func init() {
	tm := []TypeMarshaler{
		{Type(C.g_array_get_type()), marshalArray},
	}
	RegisterGValueMarshalers(tm)
}

// Array is a wrapper around GLib's GArray where each element is a GValue.
type Array struct {
	array *C.GArray
}

func wrapArray(carray *C.GArray) *Array {
	arr := &Array{array: carray}
	runtime.SetFinalizer(arr, func(a *Array) {
		C.g_array_unref(a.array)
	})
	return arr
}

// NewArray creates a new GArray from a slice of values. Each value is converted
// to a GValue and deep-copied into the array.
func NewArray(values []interface{}) (*Array, error) {
	carray := C._g_array_new_gvalues(C.guint(len(values)))

	for i, v := range values {
		gval, err := GValue(v)
		if err != nil {
			C.g_array_unref(carray)
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		C._g_array_append_gvalue(carray, gval.native())
	}

	return wrapArray(carray), nil
}

// Len returns the number of elements in the array.
func (a *Array) Len() int {
	return int(a.array.len)
}

// Index returns the Go value at position i.
func (a *Array) Index(i int) (interface{}, error) {
	if i < 0 || i >= a.Len() {
		return nil, fmt.Errorf("index %d out of bounds [0, %d)", i, a.Len())
	}
	cval := C._g_array_index_gvalue(a.array, C.guint(i))
	v := &Value{cval}
	return v.GoValue()
}

// Values returns all elements as a Go slice.
func (a *Array) Values() ([]interface{}, error) {
	n := a.Len()
	result := make([]interface{}, n)
	for i := 0; i < n; i++ {
		val, err := a.Index(i)
		if err != nil {
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		result[i] = val
	}
	return result, nil
}

// Native returns the underlying *C.GArray as an unsafe.Pointer.
func (a *Array) Native() unsafe.Pointer {
	return unsafe.Pointer(a.array)
}

// ToGValue converts the Array to a GValue of boxed type.
func (a *Array) ToGValue() (*Value, error) {
	val, err := ValueInit(Type(C.g_array_get_type()))
	if err != nil {
		return nil, err
	}
	val.SetBoxed(unsafe.Pointer(a.array))
	return val, nil
}

func marshalArray(p unsafe.Pointer) (interface{}, error) {
	carray := C.g_value_get_boxed((*C.GValue)(p))
	a := wrapArray((*C.GArray)(carray))
	C.g_array_ref(a.array)
	return a, nil
}
