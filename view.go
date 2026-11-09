package gobs

// #include <obs/obs.h>
import "C"

type View struct {
	c *C.obs_view_t
}

func ViewCreate() View {
	return View{C.obs_view_create()}
}

func (v View) Valid() bool {
	return v.c != nil
}

func (v View) Destroy() {
	C.obs_view_destroy(v.c)
}

func (v View) SetSource(channel uint32, source Source) {
	C.obs_view_set_source(v.c, C.uint32_t(channel), source.c)
}

func (v View) Add() Video {
	return Video{C.obs_view_add(v.c)}
}

func (v View) Remove() {
	C.obs_view_remove(v.c)
}
