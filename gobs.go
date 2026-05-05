package gobs

// #include <obs/obs.h>
import "C"
import (
	"unsafe"
)

func Startup(locale string) {
	cLocale := C.CString(locale)
	defer C.free(unsafe.Pointer(cLocale))

	// #cgo noescape obs_startup
	// #cgo nocallback obs_startup
	C.obs_startup(cLocale, nil, nil)
}

func Shutdown() {
	// #cgo noescape obs_shutdown
	// #cgo nocallback obs_shutdown
	C.obs_shutdown()
}

func Initialized() bool {
	// #cgo noescape obs_initialized
	// #cgo nocallback obs_initialized
	return bool(C.obs_initialized())
}

func Version() uint32 {
	// #cgo noescape obs_get_version
	// #cgo nocallback obs_get_version
	return uint32(C.obs_get_version())
}

func VersionString() string {
	// #cgo noescape obs_get_version_string
	// #cgo nocallback obs_get_version_string
	c := C.obs_get_version_string()
	return C.GoString(c)
}

func GetVideo() Video {
	// #cgo noescape obs_get_video
	// #cgo nocallback obs_get_video
	return Video{C.obs_get_video()}
}

func GetAudio() Audio {
	// #cgo noescape obs_get_audio
	// #cgo nocallback obs_get_audio
	return Audio{C.obs_get_audio()}
}

func GlobalSignalHandler() SignalHandler {
	return SignalHandler{C.obs_get_signal_handler()}
}

func EnterGraphics() {
	// #cgo noescape obs_enter_graphics
	// #cgo nocallback obs_enter_graphics
	C.obs_enter_graphics()
}

func LeaveGraphics() {
	// #cgo noescape obs_leave_graphics
	// #cgo nocallback obs_leave_graphics
	C.obs_leave_graphics()
}
