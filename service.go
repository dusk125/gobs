package gobs

// #include <obs/obs.h>
import "C"
import "unsafe"

type Service struct {
	c *C.obs_service_t
}

func (s Service) Valid() bool {
	return s.c != nil
}

func (s Service) Release() {
	// #cgo noescape obs_service_release
	// #cgo nocallback obs_service_release
	C.obs_service_release(s.c)
}

func (s Service) CanTryToConnect() bool {
	// #cgo noescape obs_service_can_try_to_connect
	// #cgo nocallback obs_service_can_try_to_connect
	return bool(C.obs_service_can_try_to_connect(s.c))
}

func ServiceCreate(id, name string, settings, hotkeys Data) Service {
	cid := C.CString(id)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cid))
	defer C.free(unsafe.Pointer(cname))

	obsSettings := settings.obs_data()
	obsHotkeys := hotkeys.obs_data()
	defer obsSettings.Release()
	defer obsHotkeys.Release()

	// #cgo noescape obs_service_create
	// #cgo nocallback obs_service_create
	service := Service{C.obs_service_create(cid, cname, obsSettings.obs_data_t, obsHotkeys.obs_data_t)}
	if service.c == nil {
		panic("failed to create service")
	}
	return service
}

func ServiceDefaults(id string) Data {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	// #cgo noescape obs_service_defaults
	// #cgo nocallback obs_service_defaults
	d := obs_data{C.obs_service_defaults(cid)}
	defer d.Release()
	return d.MapWithDefaults()
}
