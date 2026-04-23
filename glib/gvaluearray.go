package glib

// #include "glib.go.h"
//
// #pragma GCC diagnostic push
// #pragma GCC diagnostic ignored "-Wdeprecated-declarations"
//
// static GValueArray *_g_value_array_new(guint n_prealloced) {
//     return g_value_array_new(n_prealloced);
// }
//
// static GValueArray *_g_value_array_append(GValueArray *array, const GValue *value) {
//     return g_value_array_append(array, value);
// }
//
// static GValueArray *_g_value_array_copy(const GValueArray *array) {
//     return g_value_array_copy(array);
// }
//
// static void _g_value_array_free(GValueArray *array) {
//     g_value_array_free(array);
// }
//
// static GValue *_g_value_array_get_nth(GValueArray *array, guint index_) {
//     return g_value_array_get_nth(array, index_);
// }
//
// static GType _g_value_array_get_type(void) {
//     return g_value_array_get_type();
// }
//
// static guint _g_value_array_n_values(GValueArray *array) {
//     return array->n_values;
// }
//
// #pragma GCC diagnostic pop
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

var TYPE_VALUE_ARRAY Type = Type(C._g_value_array_get_type())

func init() {
	tm := []TypeMarshaler{
		{TYPE_VALUE_ARRAY, marshalValueArray},
	}
	RegisterGValueMarshalers(tm)
}

// ValueArray is a wrapper around GObject's GValueArray, a container for an
// array of GValue elements. GValueArray is deprecated in GLib since 2.32 in
// favour of GArray, but it is still the concrete type behind many object
// properties in older libraries, so bindings are required to read them.
type ValueArray struct {
	valueArray *C.GValueArray
}

func wrapValueArray(c *C.GValueArray) *ValueArray {
	va := &ValueArray{valueArray: c}
	runtime.SetFinalizer(va, func(v *ValueArray) {
		C._g_value_array_free(v.valueArray)
	})
	return va
}

// NewValueArray creates a new GValueArray from a slice of values. Each value
// is converted to a GValue and deep-copied into the array.
func NewValueArray(values []interface{}) (*ValueArray, error) {
	carray := C._g_value_array_new(C.guint(len(values)))

	for i, v := range values {
		gval, err := GValue(v)
		if err != nil {
			C._g_value_array_free(carray)
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		C._g_value_array_append(carray, gval.native())
	}

	return wrapValueArray(carray), nil
}

// Len returns the number of elements in the array.
func (a *ValueArray) Len() int {
	return int(C._g_value_array_n_values(a.valueArray))
}

// Index returns the Go value at position i.
func (a *ValueArray) Index(i int) (interface{}, error) {
	if i < 0 || i >= a.Len() {
		return nil, fmt.Errorf("index %d out of bounds [0, %d)", i, a.Len())
	}
	cval := C._g_value_array_get_nth(a.valueArray, C.guint(i))
	v := ValueFromNative(unsafe.Pointer(cval))
	return v.GoValue()
}

// Values returns all elements as a Go slice.
func (a *ValueArray) Values() ([]interface{}, error) {
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

// Copy returns a deep copy of the array.
func (a *ValueArray) Copy() *ValueArray {
	return wrapValueArray(C._g_value_array_copy(a.valueArray))
}

// Native returns the underlying *C.GValueArray as an unsafe.Pointer.
func (a *ValueArray) Native() unsafe.Pointer {
	return unsafe.Pointer(a.valueArray)
}

// ToGValue converts the ValueArray to a GValue of boxed type G_TYPE_VALUE_ARRAY.
// g_value_set_boxed deep-copies the array, so the source remains owned by the
// caller and is freed by its own finalizer.
func (a *ValueArray) ToGValue() (*Value, error) {
	val, err := ValueInit(TYPE_VALUE_ARRAY)
	if err != nil {
		return nil, err
	}
	val.SetBoxed(unsafe.Pointer(a.valueArray))
	return val, nil
}

func marshalValueArray(p unsafe.Pointer) (interface{}, error) {
	c := (*C.GValueArray)(C.g_value_get_boxed((*C.GValue)(p)))
	if c == nil {
		return (*ValueArray)(nil), nil
	}
	return wrapValueArray(C._g_value_array_copy(c)), nil
}
