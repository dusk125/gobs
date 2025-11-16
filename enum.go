package gobs

/*
#include <obs/obs.h>

bool sources_cb(void*, obs_source_t*);
*/
import "C"

import (
	"iter"
	"runtime/cgo"
	"unsafe"
)

//export sources_cb
func sources_cb(p unsafe.Pointer, s *C.obs_source_t) C.bool {
	yield := (*cgo.Handle)(p).Value().(func(Source) bool)

	return C.bool(yield(Source{s}))
}

func Sources() iter.Seq[Source] {
	// #cgo noescape obs_enum_sources
	return func(yield func(Source) bool) {
		h := cgo.NewHandle(yield)
		defer h.Delete()
		C.obs_enum_sources((*[0]byte)(C.sources_cb), unsafe.Pointer(&h))
	}
}

func SourceTypes() iter.Seq[string] {
	// #cgo noescape obs_enum_source_types
	// #cgo nocallback obs_enum_source_types
	return func(yield func(string) bool) {
		var s *CString

		for i := C.size_t(0); C.obs_enum_source_types(i, (**C.char)(unsafe.Pointer(&s))); i++ {
			if !yield(s.String()) {
				break
			}
		}
	}
}

func InputTypes() iter.Seq[string] {
	// #cgo noescape obs_enum_input_types
	// #cgo nocallback obs_enum_input_types
	return func(yield func(string) bool) {
		var s *CString
		for i := C.size_t(0); C.obs_enum_input_types(i, (**C.char)(unsafe.Pointer(&s))); i++ {
			if !yield(s.String()) {
				break
			}
		}
	}
}

func FilterTypes() iter.Seq[string] {
	// #cgo noescape obs_enum_filter_types
	// #cgo nocallback obs_enum_filter_types
	return func(yield func(string) bool) {
		var s *CString
		for i := C.size_t(0); C.obs_enum_filter_types(i, (**C.char)(unsafe.Pointer(&s))); i++ {
			if !yield(s.String()) {
				break
			}
		}
	}
}

func TransitionTypes() iter.Seq[string] {
	// #cgo noescape obs_enum_transition_types
	// #cgo nocallback obs_enum_transition_types
	return func(yield func(string) bool) {
		var s *CString
		for i := C.size_t(0); C.obs_enum_transition_types(i, (**C.char)(unsafe.Pointer(&s))); i++ {
			if !yield(s.String()) {
				break
			}
		}
	}
}

func OutputTypes() iter.Seq[string] {
	// #cgo noescape obs_enum_output_types
	// #cgo nocallback obs_enum_output_types
	return func(yield func(string) bool) {
		var s *CString
		for i := C.size_t(0); C.obs_enum_output_types(i, (**C.char)(unsafe.Pointer(&s))); i++ {
			if !yield(s.String()) {
				break
			}
		}
	}
}

func EncoderTypes() iter.Seq[string] {
	// #cgo noescape obs_enum_encoder_types
	// #cgo nocallback obs_enum_encoder_types
	return func(yield func(string) bool) {
		var s *CString
		for i := C.size_t(0); C.obs_enum_encoder_types(i, (**C.char)(unsafe.Pointer(&s))); i++ {
			if !yield(s.String()) {
				break
			}
		}
	}
}

func ServiceTypes() iter.Seq[string] {
	// #cgo noescape obs_enum_service_types
	// #cgo nocallback obs_enum_service_types
	return func(yield func(string) bool) {
		var s *CString
		for i := C.size_t(0); C.obs_enum_service_types(i, (**C.char)(unsafe.Pointer(&s))); i++ {
			if !yield(s.String()) {
				break
			}
		}
	}
}
