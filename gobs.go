package gobs

// #include <obs/obs.h>
import "C"
import (
	"unsafe"
)

func Startup(locale string, args ...string) bool {
	cLocale := C.CString(locale)
	defer C.free(unsafe.Pointer(cLocale))

	argc := C.int(len(args))
	argv := make([]*C.char, len(args))
	for i, arg := range args {
		argv[i] = C.CString(arg)
	}

	C.obs_set_cmdline_args(argc, &argv[0])

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
