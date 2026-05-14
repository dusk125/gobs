package gobs

// #include <obs/obs.h>
import "C"

type View struct {
	c *C.obs_view_t
}

func ViewCreate() View {
	// #cgo noescape obs_view_create
	// #cgo nocallback obs_view_create
	return View{C.obs_view_create()}
}

func (v View) Valid() bool {
	return v.c != nil
}

func (v View) Destroy() {
	// #cgo noescape obs_view_destroy
	// #cgo nocallback obs_view_destroy
	C.obs_view_destroy(v.c)
}

func (v View) SetSource(channel uint32, source Source) {
	// #cgo noescape obs_view_set_source
	// #cgo nocallback obs_view_set_source
	C.obs_view_set_source(v.c, C.uint32_t(channel), source.c)
}

func (v View) Add() Video {
	// #cgo noescape obs_view_add
	// #cgo nocallback obs_view_add
	return Video{C.obs_view_add(v.c)}
}

func (v View) Remove() {
	// #cgo noescape obs_view_remove
	// #cgo nocallback obs_view_remove
	C.obs_view_remove(v.c)
}
