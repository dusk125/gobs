package graphics

// #include <obs/graphics/graphics.h>
import "C"

func StageTexture(dst StageSurface, src Texture) {
	// #cgo noescape gs_stage_texture
	// #cgo nocallback gs_stage_texture
	C.gs_stage_texture(dst.c, src.c)
}

func StageSurfaceCreate(width, height uint32, format ColorFormat) StageSurface {
	// #cgo noescape gs_stagesurface_create
	// #cgo nocallback gs_stagesurface_create
	return StageSurface{C.gs_stagesurface_create(C.uint32_t(width), C.uint32_t(height), C.enum_gs_color_format(format))}
}

type StageSurface struct {
	c *C.gs_stagesurf_t
}

func (s StageSurface) Destroy() {
	// #cgo noescape gs_stagesurface_destroy
	// #cgo nocallback gs_stagesurface_destroy
	C.gs_stagesurface_destroy(s.c)
}

func (s StageSurface) Map() (*uint8, uint32, bool) {
	var data *C.uint8_t
	var lineSize C.uint32_t
	// #cgo noescape gs_stagesurface_map
	// #cgo nocallback gs_stagesurface_map
	ok := bool(C.gs_stagesurface_map(s.c, &data, &lineSize))
	if !ok {
		return nil, 0, ok
	}

	return (*uint8)(data), uint32(lineSize), true
}

func (s StageSurface) Unmap() {
	// #cgo noescape gs_stagesurface_unmap
	// #cgo nocallback gs_stagesurface_unmap
	C.gs_stagesurface_unmap(s.c)
}
