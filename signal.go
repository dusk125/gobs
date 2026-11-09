package gobs

/*
#include <stdlib.h>
#include <obs/callback/signal.h>

extern void signal_cb(void *, calldata_t*);

static void gobs_signal_connect(signal_handler_t *h, const char *sig, uintptr_t data) {
	signal_handler_connect(h, sig, signal_cb, (void *)data);
}

static void gobs_signal_disconnect(signal_handler_t *h, const char *sig, uintptr_t data) {
	signal_handler_disconnect(h, sig, signal_cb, (void *)data);
}
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type Callback func(data any, cd CallData)

type Signal interface {
	isSignal()
}

type SignalHandler struct {
	c *C.signal_handler_t
}

func (sh SignalHandler) Valid() bool {
	return sh.c != nil
}

type internal_calldata struct {
	callback Callback
	data     any
}

//export signal_cb
func signal_cb(data unsafe.Pointer, cd *C.calldata_t) {
	h := cgo.Handle(uintptr(data))
	v := h.Value().(*internal_calldata)

	v.callback(v.data, CallData{c: cd})
}

func (sh SignalHandler) Connect(name string, callback Callback, data any) SignalConnection {
	cname := C.CString(name)
	handle := cgo.NewHandle(&internal_calldata{callback: callback, data: data})
	C.gobs_signal_connect(sh.c, cname, C.uintptr_t(handle))
	return SignalConnection{
		handler: sh,
		name:    cname,
		handle:  handle,
	}
}

type SignalConnection struct {
	handler SignalHandler
	name    *C.char
	handle  cgo.Handle
}

func (sc SignalConnection) Disconnect() {
	C.gobs_signal_disconnect(sc.handler.c, sc.name, C.uintptr_t(sc.handle))
	sc.handle.Delete()
	C.free(unsafe.Pointer(sc.name))
}
