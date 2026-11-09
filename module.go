package gobs

/*
#include <obs/obs.h>

extern void find_modules_cb(void*, struct obs_module_info2*);

static void gobs_find_modules(uintptr_t data) {
	obs_find_modules2((obs_find_module_callback2_t)find_modules_cb, (void *)data);
}
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
	if bool(C.obs_init_module(m.c)) {
		return nil
	}
	return errors.New("failed to initialize module")
}

func OpenModule(path, dataPath string) (module Module, err error) {
	cpath := fromString(path).cptr()
	cdatapath := fromString(dataPath).cptr()
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
	C.obs_add_module_path(fromString(path).cptr(), fromString(dataPath).cptr())
}

//export find_modules_cb
func find_modules_cb(p unsafe.Pointer, mi *C.struct_obs_module_info2) {
	h := cgo.Handle(uintptr(p))
	find := h.Value().(func(*ModuleInfo))
	find((*ModuleInfo)(unsafe.Pointer(mi)))
}

func FindModules() iter.Seq[*ModuleInfo] {
	return func(yield func(*ModuleInfo) bool) {
		done := false
		h := cgo.NewHandle(func(mi *ModuleInfo) {
			if !done {
				done = !yield(mi)
			}
		})
		defer h.Delete()
		C.gobs_find_modules(C.uintptr_t(h))
	}
}

func LoadAllModules() {
	C.obs_load_all_modules()
}

func PostLoadModules() {
	C.obs_post_load_modules()
}

func LogLoadedModules() {
	C.obs_log_loaded_modules()
}
