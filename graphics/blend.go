package graphics

// #include <obs/graphics/graphics.h>
import "C"

func BlendStatePush() {
	// #cgo noescape gs_blend_state_push
	// #cgo nocallback gs_blend_state_push
	C.gs_blend_state_push()
}

func BlendStatePop() {
	// #cgo noescape gs_blend_state_pop
	// #cgo nocallback gs_blend_state_pop
	C.gs_blend_state_pop()
}

type BlendType uint32

const (
	BlendTypeZero        = BlendType(C.GS_BLEND_ZERO)
	BlendTypeOne         = BlendType(C.GS_BLEND_ONE)
	BlendTypeSrccolor    = BlendType(C.GS_BLEND_SRCCOLOR)
	BlendTypeInvsrccolor = BlendType(C.GS_BLEND_INVSRCCOLOR)
	BlendTypeSrcalpha    = BlendType(C.GS_BLEND_SRCALPHA)
	BlendTypeInvsrcalpha = BlendType(C.GS_BLEND_INVSRCALPHA)
	BlendTypeDstcolor    = BlendType(C.GS_BLEND_DSTCOLOR)
	BlendTypeInvdstcolor = BlendType(C.GS_BLEND_INVDSTCOLOR)
	BlendTypeDstalpha    = BlendType(C.GS_BLEND_DSTALPHA)
	BlendTypeInvdstalpha = BlendType(C.GS_BLEND_INVDSTALPHA)
	BlendTypeSrcalphasat = BlendType(C.GS_BLEND_SRCALPHASAT)
)

func BlendFunction(src, dest BlendType) {
	// #cgo noescape gs_blend_function
	// #cgo nocallback gs_blend_function
	C.gs_blend_function(C.enum_gs_blend_type(src), C.enum_gs_blend_type(dest))
}
