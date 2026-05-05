package graphics

// #cgo pkg-config: libobs
// #include <obs/graphics/graphics.h>
import "C"
import "unsafe"

type ClearFlags uint32

const (
	ClearFlagsClearColor   = ClearFlags(C.GS_CLEAR_COLOR)   // Clears color buffer
	ClearFlagsClearDepth   = ClearFlags(C.GS_CLEAR_DEPTH)   // Clears depth buffer
	ClearFlagsClearStencil = ClearFlags(C.GS_CLEAR_STENCIL) // Clears stencil buffer
)

func Clear(flags ClearFlags, color Vec4, depth float32, stencil uint8) {
	c := (*C.struct_vec4)(unsafe.Pointer(unsafe.SliceData(color[:])))
	// #cgo noescape gs_clear
	// #cgo nocallback gs_clear
	C.gs_clear(C.uint32_t(flags), c, C.float(depth), C.uint8_t(stencil))
}

func Ortho(left, right, top, bottom, znear, zfar float32) {
	// #cgo noescape gs_ortho
	// #cgo nocallback gs_ortho
	C.gs_ortho(C.float(left), C.float(right), C.float(top), C.float(bottom), C.float(znear), C.float(zfar))
}
