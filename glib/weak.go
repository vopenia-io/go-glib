package glib

/*
#include "glib.go.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type WeakRef struct {
	GWeakRef C.GWeakRef
}

func WeakRefInit(object any) *WeakRef {
	if object == nil {
		return nil
	}
	o, ok := object.(interface{ native() *C.GObject })
	if !ok {
		return nil
	}

	obj := o.native()
	if obj == nil {
		return nil
	}
	var weakRef C.GWeakRef
	C.g_weak_ref_init(&weakRef, ((C.gpointer)(obj)))
	w := &WeakRef{GWeakRef: weakRef}
	runtime.SetFinalizer(w, func(w *WeakRef) {
		C.g_weak_ref_clear(&w.GWeakRef)
	})
	return w
}

func (weakRef *WeakRef) Clear() {
	C.g_weak_ref_clear(&weakRef.GWeakRef)
	runtime.SetFinalizer(weakRef, nil)
}

func (weakRef *WeakRef) Get() *Object {
	obj := C.g_weak_ref_get(&weakRef.GWeakRef)
	if obj == nil {
		return nil
	}
	o := wrapObject(unsafe.Pointer(obj))
	o.Unref() // g_weak_ref_get() also add a ref count
	return o
}

func (weakRef *WeakRef) Set(object *Object) {
	obj := object.native()
	if obj == nil {
		C.g_weak_ref_set(&weakRef.GWeakRef, nil)
		return
	}
	C.g_weak_ref_set(&weakRef.GWeakRef, ((C.gpointer)(obj)))
}
