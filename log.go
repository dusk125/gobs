package gobs

/*
#include <obs/obs.h>
#include <obs/util/base.h>

extern void log_handler_go_cb(int, char *, int);

static void log_handler_c(int lvl, const char *format, va_list args, void *param)
{
	char out[8192];
	int len = vsnprintf(out, sizeof(out), format, args);

	log_handler_go_cb(lvl, out, len);
}

static void set_log_handler()
{
	base_set_log_handler(log_handler_c, 0);
}
*/
import "C"

import (
	"log/slog"
	"sync"
)

func init() {
	logger := slog.Default()
	SetLogHandler(func(l slog.Level, s string) {
		switch l {
		case slog.LevelDebug:
			logger.Debug("[OBS] " + s)
		case slog.LevelError:
			logger.Error("[OBS] " + s)
		case slog.LevelInfo:
			logger.Info("[OBS] " + s)
		case slog.LevelWarn:
			logger.Warn("[OBS] " + s)
		}
	})
}

var (
	logFL sync.RWMutex
	logF  func(level slog.Level, msg string)
)

//export log_handler_go_cb
func log_handler_go_cb(lvl C.int, msg *C.char, l C.int) {
	logFL.RLock()
	defer logFL.RUnlock()

	if logF == nil {
		return
	}

	level := slog.LevelDebug
	switch lvl {
	case 100: // Error
		level = slog.LevelError
	case 200: // Warning
		level = slog.LevelWarn
	case 300: // Info
		level = slog.LevelInfo
	case 400: // Debug
		level = slog.LevelDebug
	}

	logF(level, C.GoStringN(msg, l))
}

func SetLogHandler(f func(slog.Level, string)) {
	logFL.Lock()
	defer logFL.Unlock()
	logF = f

	C.set_log_handler()
}
