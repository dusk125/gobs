package graphics

// #include <obs/graphics/graphics.h>
import "C"

func StageTexture(dst StageSurface, src Texture) {
	C.gs_stage_texture(dst.c, src.c)
}

func StageSurfaceCreate(width, height uint32, format ColorFormat) StageSurface {
	return StageSurface{C.gs_stagesurface_create(C.uint32_t(width), C.uint32_t(height), C.enum_gs_color_format(format))}
}

type StageSurface struct {
	c *C.gs_stagesurf_t
}

func (s StageSurface) Valid() bool {
	return s.c != nil
}

func (s StageSurface) Destroy() {
	C.gs_stagesurface_destroy(s.c)
}

func (s StageSurface) Map() (*uint8, uint32, bool) {
	var data *C.uint8_t
	var lineSize C.uint32_t
	ok := bool(C.gs_stagesurface_map(s.c, &data, &lineSize))
	if !ok {
		return nil, 0, ok
	}

	return (*uint8)(data), uint32(lineSize), true
}

func (s StageSurface) Unmap() {
	C.gs_stagesurface_unmap(s.c)
}
