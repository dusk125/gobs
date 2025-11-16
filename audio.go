package gobs

//#include <obs/obs.h>
import "C"
import (
	"errors"
	"unsafe"
)

type SpeakerLayout uint32

const (
	SPEAKERS_UNKNOWN = SpeakerLayout(C.SPEAKERS_UNKNOWN) /**< Unknown setting, fallback is stereo. */
	SPEAKERS_MONO    = SpeakerLayout(C.SPEAKERS_MONO)    /**< Channels: MONO */
	SPEAKERS_STEREO  = SpeakerLayout(C.SPEAKERS_STEREO)  /**< Channels: FL, FR */
	SPEAKERS_2POINT1 = SpeakerLayout(C.SPEAKERS_2POINT1) /**< Channels: FL, FR, LFE */
	SPEAKERS_4POINT0 = SpeakerLayout(C.SPEAKERS_4POINT0) /**< Channels: FL, FR, FC, RC */
	SPEAKERS_4POINT1 = SpeakerLayout(C.SPEAKERS_4POINT1) /**< Channels: FL, FR, FC, LFE, RC */
	SPEAKERS_5POINT1 = SpeakerLayout(C.SPEAKERS_5POINT1) /**< Channels: FL, FR, FC, LFE, RL, RR */
	SPEAKERS_7POINT1 = SpeakerLayout(C.SPEAKERS_7POINT1) /**< Channels: FL, FR, FC, LFE, RL, RR, SL, SR */
)

type AudioInfo struct {
	SamplesPerSec uint32
	Speakers      SpeakerLayout
}

func ResetAudio(info *AudioInfo) error {
	// #cgo noescape obs_reset_audio
	// #cgo nocallback obs_reset_audio
	if bool(C.obs_reset_audio((*C.struct_obs_audio_info)(unsafe.Pointer(info)))) {
		return nil
	}
	return errors.New("failed to reset audio")
}

type Audio struct {
	c *C.audio_t
}
