package gobs

//#include <obs/obs.h>
import "C"

import (
	"bytes"
	"encoding/json"
	"unsafe"
)

type obs_data struct {
	obs_data_t *C.obs_data_t
}

func (d obs_data) Release() {
	if d.obs_data_t == nil {
		return
	}
	// #cgo noescape obs_data_release
	// #cgo nocallback obs_data_release
	C.obs_data_release(d.obs_data_t)
}

func (d obs_data) Map() Data {
	m := make(Data)

	// #cgo noescape obs_data_release
	// #cgo nocallback obs_data_release
	v := (*CString)(C.obs_data_get_json(d.obs_data_t)).Bytes()
	if err := json.Unmarshal([]byte(v), &m); err != nil {
		panic(err)
	}

	return m
}

func (d obs_data) MapWithDefaults() Data {
	var m Data

	// #cgo noescape obs_data_get_json_with_defaults
	// #cgo nocallback obs_data_get_json_with_defaults
	v := (*CString)(C.obs_data_get_json_with_defaults(d.obs_data_t)).Bytes()
	if err := json.Unmarshal(v, &m); err != nil {
		panic(err)
	}

	return m
}

type Data map[string]any

func (d Data) obs_data() obs_data {
	if d == nil {
		return obs_data{}
	}

	b := bytes.Buffer{}
	if err := json.NewEncoder(&b).Encode(d); err != nil {
		panic(err)
	}

	if err := b.WriteByte(0); err != nil {
		panic(err)
	}

	// #cgo noescape obs_data_create_from_json
	// #cgo nocallback obs_data_create_from_json
	return obs_data{C.obs_data_create_from_json((*C.char)(unsafe.Pointer(unsafe.SliceData(b.Bytes()))))}
}
