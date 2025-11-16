package gobs

//#include <obs/obs.h>
import "C"
import "unsafe"

type Encoder struct {
	c *C.struct_obs_encoder
}

func (e Encoder) Release() {
	// #cgo noescape obs_encoder_release
	// #cgo nocallback obs_encoder_release
	C.obs_encoder_release(e.c)
}

func (e Encoder) Active() bool {
	// #cgo noescape obs_encoder_active
	// #cgo nocallback obs_encoder_active
	return bool(C.obs_encoder_active(e.c))
}

func (e Encoder) Video() Video {
	// #cgo noescape obs_encoder_video
	// #cgo nocallback obs_encoder_video
	return Video{C.obs_encoder_video(e.c)}
}

func (e Encoder) SetVideo(v Video) {
	// #cgo noescape obs_encoder_set_video
	// #cgo nocallback obs_encoder_set_video
	C.obs_encoder_set_video(e.c, v.c)
}

func (e Encoder) Audio() Audio {
	// #cgo noescape obs_encoder_audio
	// #cgo nocallback obs_encoder_audio
	return Audio{C.obs_encoder_audio(e.c)}
}

func (e Encoder) SetAudio(a Audio) {
	// #cgo noescape obs_encoder_set_audio
	// #cgo nocallback obs_encoder_set_audio
	C.obs_encoder_set_audio(e.c, a.c)
}

func EncoderDefaults(id string) Data {
	cid := fromString(id).cptr()

	// #cgo noescape obs_encoder_defaults
	// #cgo nocallback obs_encoder_defaults
	d := obs_data{C.obs_encoder_defaults(cid)}
	if d.obs_data_t == nil {
		return Data{}
	}
	defer d.Release()
	return d.MapWithDefaults()
}

func VideoEncoderCreate(id string, name string, settings Data, hotkeys Data) Encoder {
	cid := C.CString(id)
	cname := C.CString(name)

	defer C.free(unsafe.Pointer(cid))
	defer C.free(unsafe.Pointer(cname))

	obsSettings := settings.obs_data()
	obsHotkeys := settings.obs_data()

	defer obsSettings.Release()
	defer obsHotkeys.Release()

	enc := Encoder{C.obs_video_encoder_create(cid, cname, obsSettings.obs_data_t, obsHotkeys.obs_data_t)}
	if enc.c == nil {
		panic("failed to create video encoder")
	}
	return enc
}

func AudioEncoderCreate(id string, name string, mixerIndex int, settings Data, hotkeys Data) Encoder {
	cid := C.CString(id)
	cname := C.CString(name)

	defer C.free(unsafe.Pointer(cid))
	defer C.free(unsafe.Pointer(cname))

	obsSettings := settings.obs_data()
	obsHotkeys := settings.obs_data()

	defer obsSettings.Release()
	defer obsHotkeys.Release()

	// #cgo noescape obs_audio_encoder_create
	// #cgo nocallback obs_audio_encoder_create
	enc := Encoder{C.obs_audio_encoder_create(cid, cname, obsSettings.obs_data_t, C.size_t(mixerIndex), obsHotkeys.obs_data_t)}
	if enc.c == nil {
		panic("failed to create audio encoder")
	}
	return enc
}
