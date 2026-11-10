package gobs

//#include <obs/obs.h>
import "C"
import "unsafe"

type Encoder struct {
	c *C.struct_obs_encoder
}

func (e Encoder) Valid() bool {
	return e.c != nil
}

func (e Encoder) Release() {
	C.obs_encoder_release(e.c)
}

func (e Encoder) Active() bool {
	return bool(C.obs_encoder_active(e.c))
}

func (e Encoder) Video() Video {
	return Video{C.obs_encoder_video(e.c)}
}

func (e Encoder) SetVideo(v Video) {
	C.obs_encoder_set_video(e.c, v.c)
}

func (e Encoder) Audio() Audio {
	return Audio{C.obs_encoder_audio(e.c)}
}

func (e Encoder) SetAudio(a Audio) {
	C.obs_encoder_set_audio(e.c, a.c)
}

func EncoderDefaults(id string) Data {
	cid := fromString(id).cptr()

	d := obs_data{C.obs_encoder_defaults(cid)}
	if d.obs_data_t == nil {
		return Data{}
	}
	defer d.Release()
	return d.MapWithDefaults()
}

func VideoEncoderCreate(id string, name string, settings Data, hotkeys Data) (Encoder, Releaser) {
	cid := C.CString(id)
	cname := C.CString(name)

	defer C.free(unsafe.Pointer(cid))
	defer C.free(unsafe.Pointer(cname))

	obsSettings := settings.obs_data()
	obsHotkeys := hotkeys.obs_data()

	defer obsSettings.Release()
	defer obsHotkeys.Release()

	enc := Encoder{C.obs_video_encoder_create(cid, cname, obsSettings.obs_data_t, obsHotkeys.obs_data_t)}
	if enc.c == nil {
		panic("failed to create video encoder")
	}
	return enc, enc.Release
}

func AudioEncoderCreate(id string, name string, mixerIndex int, settings Data, hotkeys Data) (Encoder, Releaser) {
	cid := C.CString(id)
	cname := C.CString(name)

	defer C.free(unsafe.Pointer(cid))
	defer C.free(unsafe.Pointer(cname))

	obsSettings := settings.obs_data()
	obsHotkeys := hotkeys.obs_data()

	defer obsSettings.Release()
	defer obsHotkeys.Release()

	enc := Encoder{C.obs_audio_encoder_create(cid, cname, obsSettings.obs_data_t, C.size_t(mixerIndex), obsHotkeys.obs_data_t)}
	if enc.c == nil {
		panic("failed to create audio encoder")
	}
	return enc, enc.Release
}
