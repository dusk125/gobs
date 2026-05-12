package gobs

// #include <obs/obs.h>
import "C"

import (
	"context"
	"errors"
)

type Output struct {
	c *C.obs_output_t
}

func (o Output) Valid() bool {
	return o.c != nil
}

func (o Output) Release() {
	// #cgo noescape obs_output_release
	// #cgo nocallback obs_output_release
	C.obs_output_release(o.c)
}

func (o Output) Start() error {
	// #cgo noescape obs_output_start
	// #cgo nocallback obs_output_start
	if bool(C.obs_output_start(o.c)) {
		return nil
	}
	return errors.New(o.Error())
}

func (o Output) Stop() {
	// #cgo noescape obs_output_stop
	// #cgo nocallback obs_output_stop
	C.obs_output_stop(o.c)
}

func (o Output) ForceStop() {
	// #cgo noescape obs_output_force_stop
	// #cgo nocallback obs_output_force_stop
	C.obs_output_force_stop(o.c)
}

func (o Output) Active() bool {
	// #cgo noescape obs_output_active
	// #cgo nocallback obs_output_active
	return bool(C.obs_output_active(o.c))
}

func (o Output) SetVideoEncoder(enc Encoder) {
	// #cgo noescape obs_output_set_video_encoder
	// #cgo nocallback obs_output_set_video_encoder
	C.obs_output_set_video_encoder(o.c, enc.c)
}

func (o Output) SetAudioEncoder(enc Encoder, idx int) {
	// #cgo noescape obs_output_set_audio_encoder
	// #cgo nocallback obs_output_set_audio_encoder
	C.obs_output_set_audio_encoder(o.c, enc.c, C.size_t(idx))
}

func (o Output) SetService(serv Service) {
	// #cgo noescape obs_output_set_service
	// #cgo nocallback obs_output_set_service
	C.obs_output_set_service(o.c, serv.c)
}

func (o Output) SetPreferredSize(width, height uint32) {
	// #cgo noescape obs_output_set_preferred_size
	// #cgo nocallback obs_output_set_preferred_size
	C.obs_output_set_preferred_size(o.c, C.uint32_t(width), C.uint32_t(height))
}

func (o Output) Error() string {
	// #cgo noescape obs_output_get_last_error
	// #cgo nocallback obs_output_get_last_error
	return C.GoString(C.obs_output_get_last_error(o.c))
}

func (o Output) SignalHandler() SignalHandler {
	// #cgo noescape obs_output_get_signal_handler
	// #cgo nocallback obs_output_get_signal_handler
	return SignalHandler{C.obs_output_get_signal_handler(o.c)}
}

func (o Output) TotalBytes() uint64 {
	// #cgo noescape obs_output_get_total_bytes
	// #cgo nocallback obs_output_get_total_bytes
	return uint64(C.obs_output_get_total_bytes(o.c))
}

func (o Output) FramesDropped() int32 {
	// #cgo noescape obs_output_get_frames_dropped
	// #cgo nocallback obs_output_get_frames_dropped
	return int32(C.obs_output_get_frames_dropped(o.c))
}

func (o Output) TotalFrames() int32 {
	// #cgo noescape obs_output_get_total_frames
	// #cgo nocallback obs_output_get_total_frames
	return int32(C.obs_output_get_total_frames(o.c))
}

func (o Output) Reconnecting() bool {
	// #cgo noescape obs_output_reconnecting
	// #cgo nocallback obs_output_reconnecting
	return bool(C.obs_output_reconnecting(o.c))
}

func (o Output) Events(ctx context.Context) <-chan SignalOutput {
	ch := make(chan SignalOutput, 8)

	sh := o.SignalHandler()
	connections := []SignalConnection{
		sh.Connect("start", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputStart{
				Output: Output{(*C.obs_output_t)(cd.Ptr("output"))},
			}
		}, ch),
		sh.Connect("stop", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputStop{
				Output: Output{(*C.obs_output_t)(cd.Ptr("output"))},
				Code:   OutputStopCode(cd.Int("code")),
			}
		}, ch),
		sh.Connect("pause", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputPause{
				Output: Output{(*C.obs_output_t)(cd.Ptr("output"))},
			}
		}, ch),
		sh.Connect("unpause", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputUnpause{
				Output: Output{(*C.obs_output_t)(cd.Ptr("unpause"))},
			}
		}, ch),
		sh.Connect("starting", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputStarting{
				Output: Output{(*C.obs_output_t)(cd.Ptr("starting"))},
			}
		}, ch),
		sh.Connect("stopping", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputStopping{
				Output: Output{(*C.obs_output_t)(cd.Ptr("stopping"))},
			}
		}, ch),
		sh.Connect("activate", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputActivate{
				Output: Output{(*C.obs_output_t)(cd.Ptr("activate"))},
			}
		}, ch),
		sh.Connect("deactivate", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputDeactivate{
				Output: Output{(*C.obs_output_t)(cd.Ptr("deactivate"))},
			}
		}, ch),
		sh.Connect("reconnect", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputReconnect{
				Output: Output{(*C.obs_output_t)(cd.Ptr("reconnect"))},
			}
		}, ch),
		sh.Connect("reconnect_success", func(data any, cd CallData) {
			data.(chan SignalOutput) <- SignalOutputReconnectSuccess{
				Output: Output{(*C.obs_output_t)(cd.Ptr("reconnect_success"))},
			}
		}, ch),
	}

	go func() {
		<-ctx.Done()
		for _, connection := range connections {
			connection.Disconnect()
		}
		close(ch)
	}()

	return ch
}

func OutputDefaults(id string) Data {
	cid := fromString(id)

	// #cgo noescape obs_output_defaults
	// #cgo nocallback obs_output_defaults
	d := obs_data{C.obs_output_defaults(cid.cptr())}
	defer d.Release()
	return d.MapWithDefaults()
}

func SetOutputSource(channel uint32, source Source) {
	// #cgo noescape obs_set_output_source
	// #cgo nocallback obs_set_output_source
	C.obs_set_output_source(C.uint32_t(channel), source.c)
}

// Gets the primary output source for a channel and increments the reference counter for that source.
// Use Source.Release() to release.
func OutputSource(channel uint32) Source {
	// #cgo noescape obs_get_output_source
	// #cgo nocallback obs_get_output_source
	return Source{C.obs_get_output_source(C.uint32_t(channel))}
}

func OutputCreate(id string, name string, settings Data, hotkeys Data) Output {
	cid := fromString(id).cptr()
	cname := fromString(name).cptr()

	obsSettings := settings.obs_data()
	obsHotkeys := hotkeys.obs_data()
	defer obsSettings.Release()
	defer obsHotkeys.Release()

	// #cgo noescape obs_output_create
	// #cgo nocallback obs_output_create
	out := Output{C.obs_output_create(cid, cname, obsSettings.obs_data_t, obsHotkeys.obs_data_t)}
	if out.c == nil {
		panic("failed to create output")
	}
	return out
}

type SignalOutput interface {
	Signal
	isSignalOutput()
}

type SignalOutputStart struct {
	Output Output
}

func (SignalOutputStart) isSignal()       {}
func (SignalOutputStart) isSignalOutput() {}

type OutputStopCode int64

const (
	OUTPUT_STOP_SUCCESS        OutputStopCode = C.OBS_OUTPUT_SUCCESS        // Successfully stopped
	OUTPUT_STOP_BAD_PATH       OutputStopCode = C.OBS_OUTPUT_BAD_PATH       // The specified path was invalid
	OUTPUT_STOP_CONNECT_FAILED OutputStopCode = C.OBS_OUTPUT_CONNECT_FAILED // Failed to connect to a server
	OUTPUT_STOP_INVALID_STREAM OutputStopCode = C.OBS_OUTPUT_INVALID_STREAM // Invalid stream path
	OUTPUT_STOP_ERROR          OutputStopCode = C.OBS_OUTPUT_ERROR          // Generic error
	OUTPUT_STOP_DISCONNECTED   OutputStopCode = C.OBS_OUTPUT_DISCONNECTED   // Unexpectedly disconnected
	OUTPUT_STOP_UNSUPPORTED    OutputStopCode = C.OBS_OUTPUT_UNSUPPORTED    // The settings, video/audio format, or codecs are unsupported by this output
	OUTPUT_STOP_NO_SPACE       OutputStopCode = C.OBS_OUTPUT_NO_SPACE       // Ran out of disk space
	OUTPUT_STOP_ENCODE_ERROR   OutputStopCode = C.OBS_OUTPUT_ENCODE_ERROR   // Encoder error
)

type SignalOutputStop struct {
	Output Output
	Code   OutputStopCode
}

func (SignalOutputStop) isSignal()       {}
func (SignalOutputStop) isSignalOutput() {}

type SignalOutputPause struct {
	Output Output
}

func (SignalOutputPause) isSignal()       {}
func (SignalOutputPause) isSignalOutput() {}

type SignalOutputUnpause struct {
	Output Output
}

func (SignalOutputUnpause) isSignal()       {}
func (SignalOutputUnpause) isSignalOutput() {}

type SignalOutputStarting struct {
	Output Output
}

func (SignalOutputStarting) isSignal()       {}
func (SignalOutputStarting) isSignalOutput() {}

type SignalOutputStopping struct {
	Output Output
}

func (SignalOutputStopping) isSignal()       {}
func (SignalOutputStopping) isSignalOutput() {}

type SignalOutputActivate struct {
	Output Output
}

func (SignalOutputActivate) isSignal()       {}
func (SignalOutputActivate) isSignalOutput() {}

type SignalOutputDeactivate struct {
	Output Output
}

func (SignalOutputDeactivate) isSignal()       {}
func (SignalOutputDeactivate) isSignalOutput() {}

type SignalOutputReconnect struct {
	Output Output
}

func (SignalOutputReconnect) isSignal()       {}
func (SignalOutputReconnect) isSignalOutput() {}

type SignalOutputReconnectSuccess struct {
	Output Output
}

func (SignalOutputReconnectSuccess) isSignal()       {}
func (SignalOutputReconnectSuccess) isSignalOutput() {}
