package gobs

/*
#include <obs/callback/proc.h>

extern void proc_handler_cb(void *, calldata_t*);
*/
import "C"

type ProcHandler struct {
	c *C.proc_handler_t
}

func (ph ProcHandler) Valid() bool {
	return ph.c != nil
}

func (ph ProcHandler) Call(name string, params CallData) bool {
	return bool(C.proc_handler_call(ph.c, fromString(name).cptr(), params.c))
}
