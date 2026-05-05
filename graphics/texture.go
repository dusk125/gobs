package graphics

// #include <obs/graphics/graphics.h>
import "C"

type Texture struct {
	c *C.gs_texture_t
}
