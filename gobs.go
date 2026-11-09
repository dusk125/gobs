package gobs

// #include <obs/obs.h>
import "C"
import (
	"unsafe"
)

func Startup(locale string, args ...string) bool {
	cLocale := C.CString(locale)
	defer C.free(unsafe.Pointer(cLocale))

	argv := make([]*C.char, len(args)+1)
	for i, arg := range args {
		argv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(argv[i]))
	}

	C.obs_set_cmdline_args(C.int(len(args)), unsafe.SliceData(argv))

	return bool(C.obs_startup(cLocale, nil, nil))
}

func Shutdown() {
	C.obs_shutdown()
}

func Initialized() bool {
	return bool(C.obs_initialized())
}

func Version() uint32 {
	return uint32(C.obs_get_version())
}

func VersionString() string {
	c := C.obs_get_version_string()
	return C.GoString(c)
}

func GetVideo() Video {
	return Video{C.obs_get_video()}
}

func GetAudio() Audio {
	return Audio{C.obs_get_audio()}
}

func GlobalSignalHandler() SignalHandler {
	return SignalHandler{C.obs_get_signal_handler()}
}

func EnterGraphics() {
	C.obs_enter_graphics()
}

func LeaveGraphics() {
	C.obs_leave_graphics()
}
