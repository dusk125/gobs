package gobs

//#include <obs/obs.h>
import "C"
import (
	"fmt"
	"reflect"
	"runtime/cgo"
	"unsafe"
)

type CallData struct {
	c       *C.calldata_t
	strings []*C.char
	handles []cgo.Handle
}

func (cd *CallData) init() {
	if cd.c == nil {
		cd.c = &C.calldata_t{}
	}
}

func (cd *CallData) Free() {
	for _, str := range cd.strings {
		C.free(unsafe.Pointer(str))
	}
	for _, handle := range cd.handles {
		handle.Delete()
	}
}

func (cd *CallData) Set(name string, val any) {
	cd.init()

	cname := C.CString(name)
	cd.strings = append(cd.strings, cname)

	switch val := val.(type) {
	case int:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case int8:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case int16:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case int32:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case int64:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case uint:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case uint8:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case uint16:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case uint32:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case uint64:
		C.calldata_set_int(cd.c, cname, C.longlong(val))
	case float32:
		C.calldata_set_float(cd.c, cname, C.double(val))
	case float64:
		C.calldata_set_float(cd.c, cname, C.double(val))
	case bool:
		C.calldata_set_bool(cd.c, cname, C.bool(val))
	case string:
		cval := C.CString(val)
		cd.strings = append(cd.strings, cval)
		C.calldata_set_string(cd.c, cname, cval)
	default:
		v := reflect.ValueOf(val)
		if v.Kind() != reflect.Struct {
			panic(fmt.Errorf("invalid calldata type: %T", val))
		}
		c := v.FieldByName("c")
		if c.IsZero() {
			panic(fmt.Errorf("invalid calldata type, expect a c-like structure: %T", val))
		}
		C.calldata_set_ptr(cd.c, cname, c.UnsafePointer())
	}
}

func (cd CallData) Int(name string) int64 {
	cname := fromString(name).cptr()
	return int64(C.calldata_int(cd.c, cname))
}

func (cd CallData) Float(name string) float64 {
	cname := fromString(name).cptr()
	return float64(C.calldata_float(cd.c, cname))
}

func (cd CallData) Bool(name string) bool {
	cname := fromString(name).cptr()
	return bool(C.calldata_bool(cd.c, cname))
}

func (cd CallData) Ptr(name string) unsafe.Pointer {
	cname := fromString(name).cptr()
	return C.calldata_ptr(cd.c, cname)
}

func (cd CallData) String(name string) string {
	cname := fromString(name).cptr()
	return (*CString)(C.calldata_ptr(cd.c, cname)).String()
}
