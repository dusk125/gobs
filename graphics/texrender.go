package graphics

// #include <obs/graphics/graphics.h>
import "C"

type TexRender struct {
	c *C.gs_texrender_t
}

type ColorFormat uint32

const (
	ColorFormatUNKNOWN     = ColorFormat(C.GS_UNKNOWN)     // Unknown format
	ColorFormatA8          = ColorFormat(C.GS_A8)          // 8 bit alpha channel only
	ColorFormatR8          = ColorFormat(C.GS_R8)          // 8 bit red channel only
	ColorFormatRGBA        = ColorFormat(C.GS_RGBA)        // RGBA, 8 bits per channel
	ColorFormatBGRX        = ColorFormat(C.GS_BGRX)        // BGRX, 8 bits per channel
	ColorFormatBGRA        = ColorFormat(C.GS_BGRA)        // BGRA, 8 bits per channel
	ColorFormatR10G10B10A2 = ColorFormat(C.GS_R10G10B10A2) // RGBA, 10 bits per channel except alpha, which is 2 bits
	ColorFormatRGBA16      = ColorFormat(C.GS_RGBA16)      // RGBA, 16 bits per channel
	ColorFormatR16         = ColorFormat(C.GS_R16)         // 16 bit red channel only
	ColorFormatRGBA16F     = ColorFormat(C.GS_RGBA16F)     // RGBA, 16 bit floating point per channel
	ColorFormatRGBA32F     = ColorFormat(C.GS_RGBA32F)     // RGBA, 32 bit floating point per channel
	ColorFormatRG16F       = ColorFormat(C.GS_RG16F)       // 16 bit floating point red and green channels only
	ColorFormatRG32F       = ColorFormat(C.GS_RG32F)       // 32 bit floating point red and green channels only
	ColorFormatR16F        = ColorFormat(C.GS_R16F)        // 16 bit floating point red channel only
	ColorFormatR32F        = ColorFormat(C.GS_R32F)        // 32 bit floating point red channel only
	ColorFormatDXT1        = ColorFormat(C.GS_DXT1)        // Compressed DXT1
	ColorFormatDXT3        = ColorFormat(C.GS_DXT3)        // Compressed DXT3
	ColorFormatDXT5        = ColorFormat(C.GS_DXT5)        // Compressed DXT5
	ColorFormatRGBA_UNORM  = ColorFormat(C.GS_RGBA_UNORM)  // RGBA, 8 bits per channel, no SRGB aliasing
	ColorFormatBGRX_UNORM  = ColorFormat(C.GS_BGRX_UNORM)  // BGRX, 8 bits per channel, no SRGB aliasing
	ColorFormatBGRA_UNORM  = ColorFormat(C.GS_BGRA_UNORM)  // BGRA, 8 bits per channel, no SRGB aliasing
	ColorFormatRG16        = ColorFormat(C.GS_RG16)        // RG, 16 bits per channel
)

type ZStencilFormat uint32

const (
	ZStencilNONE      = ZStencilFormat(C.GS_ZS_NONE)    // No Z-stencil buffer
	ZStencil16        = ZStencilFormat(C.GS_Z16)        // 16 bit Z buffer
	ZStencil24_S8     = ZStencilFormat(C.GS_Z24_S8)     // 24 bit Z buffer, 8 bit stencil
	ZStencil32F       = ZStencilFormat(C.GS_Z32F)       // 32 bit floating point Z buffer
	ZStencil32F_S8X24 = ZStencilFormat(C.GS_Z32F_S8X24) // 32 bit floating point Z buffer, 8 bit stencil
)

func TexRenderCreate(format ColorFormat, zsformat ZStencilFormat) TexRender {
	// #cgo noescape gs_texrender_create
	// #cgo nocallback gs_texrender_create
	return TexRender{C.gs_texrender_create(C.enum_gs_color_format(format), C.enum_gs_zstencil_format(zsformat))}
}

func (t TexRender) Destroy() {
	// #cgo noescape gs_texrender_destroy
	// #cgo nocallback gs_texrender_destroy
	C.gs_texrender_destroy(t.c)
}

func (t TexRender) Reset() {
	// #cgo noescape gs_texrender_reset
	// #cgo nocallback gs_texrender_reset
	C.gs_texrender_reset(t.c)
}

func (t TexRender) Begin(width, height uint32) bool {
	// #cgo noescape gs_texrender_begin
	// #cgo nocallback gs_texrender_begin
	return bool(C.gs_texrender_begin(t.c, C.uint32_t(width), C.uint32_t(height)))
}

func (t TexRender) End() {
	// #cgo noescape gs_texrender_end
	// #cgo nocallback gs_texrender_end
	C.gs_texrender_end(t.c)
}
