package gobs

// #include <obs/obs.h>
import "C"

type MediaState uint32

const (
	MediaStateNone      = MediaState(C.OBS_MEDIA_STATE_NONE)
	MediaStatePlaying   = MediaState(C.OBS_MEDIA_STATE_PLAYING)
	MediaStateOpening   = MediaState(C.OBS_MEDIA_STATE_OPENING)
	MediaStateBuffering = MediaState(C.OBS_MEDIA_STATE_BUFFERING)
	MediaStatePaused    = MediaState(C.OBS_MEDIA_STATE_PAUSED)
	MediaStateStopped   = MediaState(C.OBS_MEDIA_STATE_STOPPED)
	MediaStateEnded     = MediaState(C.OBS_MEDIA_STATE_ENDED)
	MediaStateError     = MediaState(C.OBS_MEDIA_STATE_ERROR)
)

type MediaSource struct {
	Source
}

func (s MediaSource) Restart() {
	C.obs_source_media_restart(s.c)
}

func (s MediaSource) Stop() {
	C.obs_source_media_stop(s.c)
}

func (s MediaSource) MediaTime() int64 {
	return int64(C.obs_source_media_get_time(s.c))
}

func (s MediaSource) MediaDuration() int64 {
	return int64(C.obs_source_media_get_duration(s.c))
}

func (s MediaSource) State() MediaState {
	return MediaState(C.obs_source_media_get_state(s.c))
}
