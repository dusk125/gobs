package gobs

/*
#include <obs/obs.h>

extern void findModulesCB(void*, struct obs_module_info2*);
*/
import "C"

import (
	"errors"
	"iter"
	"runtime/cgo"
	"unsafe"
)

type ModuleInfo struct {
	binPath  *CString
	dataPath *CString
	name     *CString
}

func (m ModuleInfo) BinPath() string {
	return m.binPath.String()
}

func (m ModuleInfo) DataPath() string {
	return m.dataPath.String()
}

func (m ModuleInfo) Name() string {
	return m.name.String()
}

type Module struct {
	c *C.obs_module_t
}

func (m Module) Valid() bool {
	return m.c != nil
}

func (m Module) Init() error {
	// #cgo noescape obs_init_module
	// #cgo nocallback obs_init_module
	if bool(C.obs_init_module(m.c)) {
		return nil
	}
	return errors.New("failed to initialize module")
}

func OpenModule(path, dataPath string) (module Module, err error) {
	cpath := fromString(path).cptr()
	cdatapath := fromString(dataPath).cptr()
	// #cgo noescape obs_open_module
	// #cgo nocallback obs_open_module
	switch C.obs_open_module(&module.c, cpath, cdatapath) {
	case 0:
		return
	case -1:
		return Module{}, errors.New("A generic error occurred")
	case -2:
		return Module{}, errors.New("The module was not found")
	case -3:
		return Module{}, errors.New("Required exports are missing")
	case -4:
		return Module{}, errors.New("Incompatible version")
	case -5:
		return Module{}, errors.New("Skipped by hardcoded rules (e.g. obsolete obs-browser macOS plugin)")
	default:
		return Module{}, errors.New("Unknown error")
	}
}

func AddModulePath(path, dataPath string) {
	// #cgo noescape obs_add_module_path
	// #cgo nocallback obs_add_module_path
	C.obs_add_module_path(fromString(path).cptr(), fromString(dataPath).cptr())
}

//export findModulesCB
func findModulesCB(p unsafe.Pointer, mi *C.struct_obs_module_info2) {
	h := (*cgo.Handle)(p)
	yield := h.Value().(func(*ModuleInfo) bool)
	yield((*ModuleInfo)(unsafe.Pointer(mi)))
}

func FindModules() iter.Seq[*ModuleInfo] {
	return func(yield func(*ModuleInfo) bool) {
		h := cgo.NewHandle(yield)
		defer h.Delete()
		C.obs_find_modules2((*[0]byte)(C.findModulesCB), unsafe.Pointer(&h))
	}
}

func LoadAllModules() {
	// #cgo noescape obs_load_all_modules
	// #cgo nocallback obs_load_all_modules
	C.obs_load_all_modules()
}

func PostLoadModules() {
	// #cgo noescape obs_post_load_modules
	// #cgo nocallback obs_post_load_modules
	C.obs_post_load_modules()
}

func LogLoadedModules() {
	// #cgo noescape obs_log_loaded_modules
	// #cgo nocallback obs_log_loaded_modules
	C.obs_log_loaded_modules()
}
