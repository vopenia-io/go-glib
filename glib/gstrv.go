package glib

// #include "glib.go.h"
//
// static guint _g_strv_length(gchar **strv) {
//     return g_strv_length(strv);
// }
//
// static gchar **_g_strv_new(guint len) {
//     return (gchar **)g_new0(gchar *, len + 1);
// }
//
// static void _g_strv_set(gchar **strv, guint i, const gchar *s) {
//     strv[i] = g_strdup(s);
// }
//
// static gchar *_g_strv_index(gchar **strv, guint i) {
//     return strv[i];
// }
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

var TYPE_STRV Type = Type(C.G_TYPE_STRV) // is function g_strv_get_type() inside macro, can't be const in Go

func init() {
	tm := []TypeMarshaler{
		{Type(C.g_strv_get_type()), marshalStrv},
	}
	RegisterGValueMarshalers(tm)
}

// Strv is a wrapper around a null-terminated gchar** (GStrv).
type Strv struct {
	strv **C.gchar
}

func wrapStrv(cstrv **C.gchar) *Strv {
	s := &Strv{strv: cstrv}
	runtime.SetFinalizer(s, func(s *Strv) {
		C.g_strfreev(s.strv)
	})
	return s
}

// NewStrv creates a new Strv from a Go string slice.
func NewStrv(values []string) *Strv {
	cstrv := C._g_strv_new(C.guint(len(values)))
	for i, v := range values {
		cstr := C.CString(v)
		C._g_strv_set(cstrv, C.guint(i), (*C.gchar)(cstr))
		C.free(unsafe.Pointer(cstr))
	}
	return wrapStrv(cstrv)
}

// Len returns the number of strings in the Strv.
func (s *Strv) Len() int {
	return int(C._g_strv_length(s.strv))
}

// Index returns the string at position i.
func (s *Strv) Index(i int) (string, error) {
	if i < 0 || i >= s.Len() {
		return "", fmt.Errorf("index %d out of bounds [0, %d)", i, s.Len())
	}
	return C.GoString((*C.char)(C._g_strv_index(s.strv, C.guint(i)))), nil
}

// Strings returns all elements as a Go string slice.
func (s *Strv) Strings() []string {
	n := s.Len()
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = C.GoString((*C.char)(C._g_strv_index(s.strv, C.guint(i))))
	}
	return result
}

// Native returns the underlying **C.gchar as an unsafe.Pointer.
func (s *Strv) Native() unsafe.Pointer {
	return unsafe.Pointer(s.strv)
}

// ToGValue converts the Strv to a GValue of boxed type.
func (s *Strv) ToGValue() (*Value, error) {
	val, err := ValueInit(Type(C.g_strv_get_type()))
	if err != nil {
		return nil, err
	}
	val.SetBoxed(unsafe.Pointer(s.strv))
	return val, nil
}

func marshalStrv(p unsafe.Pointer) (interface{}, error) {
	cstrv := C.g_value_get_boxed((*C.GValue)(p))
	dup := C.g_strdupv((**C.gchar)(cstrv))
	return wrapStrv(dup), nil
}
