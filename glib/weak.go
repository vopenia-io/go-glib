package glib

/*
#include "glib.go.h"
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type WeakRef struct {
	GWeakRef *C.GWeakRef
}

func WeakRefInit(object any) *WeakRef {
	if object == nil {
		return nil
	}
	if !IsGObject(object) {
		panic(fmt.Sprintf("object of type %T is not a GObject", object))
	}
	o := object.(interface{ native() *C.GObject })
	obj := o.native()
	var weakRef *C.GWeakRef
	weakRef = (*C.GWeakRef)(C.malloc(C.sizeof_GWeakRef))
	C.g_weak_ref_init(weakRef, ((C.gpointer)(obj)))
	w := &WeakRef{GWeakRef: weakRef}
	runtime.SetFinalizer(w, func(w *WeakRef) {
		C.g_weak_ref_clear(w.GWeakRef)
		C.free(unsafe.Pointer(w.GWeakRef))
	})
	return w
}

func (weakRef *WeakRef) Get() *Object {
	obj := C.g_weak_ref_get(weakRef.GWeakRef)
	if obj == nil {
		return nil
	}
	o := wrapObjectClean(unsafe.Pointer(obj))
	runtime.SetFinalizer(o, (*Object).Unref)
	return o
}

func (weakRef *WeakRef) Set(object any) {
	if object == nil {
		return
	}
	if !IsGObject(object) {
		panic(fmt.Sprintf("object of type %T is not a GObject", object))
	}
	o := object.(interface{ native() *C.GObject })
	obj := o.native()
	C.g_weak_ref_set(weakRef.GWeakRef, ((C.gpointer)(obj)))
}
