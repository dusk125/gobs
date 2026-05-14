package gobs

/*
#include <obs/obs.h>

extern void main_render_cb(void *, uint32_t, uint32_t);
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type RenderCallback func(data any, cx, cy uint32)

type internal_render_calldata struct {
	callback RenderCallback
	data     any
}

//export main_render_cb
func main_render_cb(data unsafe.Pointer, cx C.uint32_t, cy C.uint32_t) {
	h := *(*cgo.Handle)(data)
	v := h.Value().(*internal_render_calldata)

	v.callback(v.data, uint32(cx), uint32(cy))
}

type RenderConnection struct {
	handle *cgo.Handle
}

func AddMainRenderCallback(callback RenderCallback, data any) RenderConnection {
	handle := cgo.NewHandle(&internal_render_calldata{callback: callback, data: data})
	C.obs_add_main_render_callback((*[0]byte)(C.main_render_cb), unsafe.Pointer(&handle))
	return RenderConnection{handle: &handle}
}

func (rc RenderConnection) Remove() {
	C.obs_remove_main_render_callback((*[0]byte)(C.main_render_cb), unsafe.Pointer(rc.handle))
	(*rc.handle).Delete()
}
