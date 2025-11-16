package gobs

import "C"

import "unsafe"

type CString C.char

func (c *CString) cptr() *C.char {
	return (*C.char)(c)
}

func (s *CString) Len() (l int) {
	for v := unsafe.Pointer(s); *(*byte)(v) != 0; l++ {
		v = unsafe.Add(v, 1)
	}
	return
}

// String returns a Go copy of the string
func (s *CString) String() string {
	return string(s.Bytes())
}

func (s *CString) Bytes() []byte {
	b := make([]byte, s.Len())
	copy(b, unsafe.Slice((*byte)(unsafe.Pointer(s)), len(b)))
	return b
}

func fromString(s string) *CString {
	return (*CString)(unsafe.Pointer(unsafe.StringData(s + "\x00")))
}
