package gobs

// #include <obs/obs.h>
import "C"
import "context"

type Source struct {
	c *C.obs_source_t
}

func SourceDefaults(id string) Data {
	cid := fromString(id).cptr()

	// #cgo noescape obs_get_source_defaults
	// #cgo nocallback obs_get_source_defaults
	data := obs_data{C.obs_get_source_defaults(cid)}
	defer data.Release()
	return data.MapWithDefaults()
}

func SourceCreate(id string, name string, settings Data, hotkey Data) Source {
	var (
		obsSettings obs_data
		obsHotkey   obs_data
	)

	if settings != nil {
		obsSettings = settings.obs_data()
		defer obsSettings.Release()
	}
	if hotkey != nil {
		obsHotkey = hotkey.obs_data()
		defer obsHotkey.Release()
	}

	// #cgo noescape obs_source_create
	// #cgo nocallback obs_source_create
	source := Source{C.obs_source_create(fromString(id).cptr(), fromString(name).cptr(), obsSettings.obs_data_t, obsHotkey.obs_data_t)}
	if source.c == nil {
		panic("failed to create source")
	}
	return source
}

func (s Source) Release() {
	// #cgo noescape obs_source_release
	// #cgo nocallback obs_source_release
	C.obs_source_release(s.c)
}

func (s Source) Update(settings Data) {
	d := settings.obs_data()
	defer d.Release()

	// #cgo noescape obs_source_update
	// #cgo nocallback obs_source_update
	C.obs_source_update(s.c, d.obs_data_t)
}

func (s Source) Name() string {
	// #cgo noescape obs_source_get_name
	// #cgo nocallback obs_source_get_name
	return C.GoString(C.obs_source_get_name(s.c))
}

func (s Source) ProcHandler() ProcHandler {
	return ProcHandler{C.obs_source_get_proc_handler(s.c)}
}

func (s Source) SignalHandler() SignalHandler {
	return SignalHandler{C.obs_source_get_signal_handler(s.c)}
}

func (s Source) Events(ctx context.Context) <-chan SignalSource {
	ch := make(chan SignalSource, 8)

	sh := s.SignalHandler()
	connections := []SignalConnection{
		sh.Connect("destroy", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceDestroy{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("remove", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceRemove{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("update", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceUpdate{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("save", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceSave{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("load", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceLoad{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("activate", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceActivate{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("deactivate", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceDeactivate{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("show", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceShow{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("hide", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceHide{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("mute", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMute{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Muted:  cd.Bool("muted"),
			}
		}, ch),
		sh.Connect("push_to_mute_changed", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourcePushToMuteChanged{
				Source:  Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Enabled: cd.Bool("enabled"),
			}
		}, ch),
		sh.Connect("push_to_mute_delay", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourcePushToMuteDelay{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Delay:  cd.Int("delay"),
			}
		}, ch),
		sh.Connect("push_to_talk_changed", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourcePushToTalkChanged{
				Source:  Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Enabled: cd.Bool("enabled"),
			}
		}, ch),
		sh.Connect("push_to_talk_delay", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourcePushToTalkDelay{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Delay:  cd.Int("delay"),
			}
		}, ch),
		sh.Connect("enable", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceEnable{
				Source:  Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Enabled: cd.Bool("enabled"),
			}
		}, ch),
		sh.Connect("rename", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceRename{
				Source:   Source{(*C.obs_source_t)(cd.Ptr("source"))},
				NewName:  cd.String("new_name"),
				PrevName: cd.String("prev_name"),
			}
		}, ch),
		sh.Connect("volume", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceVolume{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Volume: cd.Float("volume"),
			}
		}, ch),
		sh.Connect("update_properties", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceUpdateProperties{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("update_flags", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceUpdateFlags{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Flags:  cd.Int("flags"),
			}
		}, ch),
		sh.Connect("audio_sync", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceAudioSync{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Offset: cd.Int("offset"),
			}
		}, ch),
		sh.Connect("audio_balance", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceAudioBalance{
				Source:  Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Balance: cd.Float("balance"),
			}
		}, ch),
		sh.Connect("audio_mixers", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceAudioMixers{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Mixers: cd.Int("mixers"),
			}
		}, ch),
		sh.Connect("audio_activate", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceAudioActivate{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("audio_deactivate", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceAudioDeactivate{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("filter_add", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceFilterAdd{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Filter: Filter{Source{(*C.obs_source_t)(cd.Ptr("filter"))}},
			}
		}, ch),
		sh.Connect("filter_remove", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceFilterRemove{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Filter: Filter{Source{(*C.obs_source_t)(cd.Ptr("filter"))}},
			}
		}, ch),
		sh.Connect("reorder_filters", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceReorderFilters{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("transition_start", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceTransitionStart{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("transition_video_stop", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceTransitionVideoStop{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("transition_stop", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceTransitionStop{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_started", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaStarted{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_ended", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaEnded{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_pause", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaPause{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_play", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaPlay{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_restart", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaRestart{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_stopped", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaStopped{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_next", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaNext{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("media_previous", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceMediaPrevious{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
			}
		}, ch),
		sh.Connect("slide_changed", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalImageSlideShowSlideChanged{
				Index: cd.Int("index"),
				Path:  cd.String("path"),
			}
		}, ch),
		sh.Connect("hooked", func(data any, cd CallData) {
			name := cd.String("title")
			if name == "" {
				name = cd.String("name")
			}
			data.(chan SignalSource) <- SignalSourceHoooked{
				Source:     Source{(*C.obs_source_t)(cd.Ptr("source"))},
				Class:      cd.String("class"),
				Title:      name,
				Executable: cd.String("executable"),
			}
		}, ch),
		sh.Connect("unhooked", func(data any, cd CallData) {
			data.(chan SignalSource) <- SignalSourceUnhooked{
				Source: Source{(*C.obs_source_t)(cd.Ptr("source"))},
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

type SignalSource interface {
	Signal
	isSignalSource()
}

type SignalSourceDestroy struct {
	Source Source
}

func (SignalSourceDestroy) isSignal()       {}
func (SignalSourceDestroy) isSignalSource() {}

type SignalSourceRemove struct {
	Source Source
}

func (SignalSourceRemove) isSignal()       {}
func (SignalSourceRemove) isSignalSource() {}

type SignalSourceUpdate struct {
	Source Source
}

func (SignalSourceUpdate) isSignal()       {}
func (SignalSourceUpdate) isSignalSource() {}

type SignalSourceSave struct {
	Source Source
}

func (SignalSourceSave) isSignal()       {}
func (SignalSourceSave) isSignalSource() {}

type SignalSourceLoad struct {
	Source Source
}

func (SignalSourceLoad) isSignal()       {}
func (SignalSourceLoad) isSignalSource() {}

type SignalSourceActivate struct {
	Source Source
}

func (SignalSourceActivate) isSignal()       {}
func (SignalSourceActivate) isSignalSource() {}

type SignalSourceDeactivate struct {
	Source Source
}

func (SignalSourceDeactivate) isSignal()       {}
func (SignalSourceDeactivate) isSignalSource() {}

type SignalSourceShow struct {
	Source Source
}

func (SignalSourceShow) isSignal()       {}
func (SignalSourceShow) isSignalSource() {}

type SignalSourceHide struct {
	Source Source
}

func (SignalSourceHide) isSignal()       {}
func (SignalSourceHide) isSignalSource() {}

type SignalSourceMute struct {
	Source Source
	Muted  bool
}

func (SignalSourceMute) isSignal()       {}
func (SignalSourceMute) isSignalSource() {}

type SignalSourcePushToMuteChanged struct {
	Source  Source
	Enabled bool
}

func (SignalSourcePushToMuteChanged) isSignal()       {}
func (SignalSourcePushToMuteChanged) isSignalSource() {}

type SignalSourcePushToMuteDelay struct {
	Source Source
	Delay  int64
}

func (SignalSourcePushToMuteDelay) isSignal()       {}
func (SignalSourcePushToMuteDelay) isSignalSource() {}

type SignalSourcePushToTalkChanged struct {
	Source  Source
	Enabled bool
}

func (SignalSourcePushToTalkChanged) isSignal()       {}
func (SignalSourcePushToTalkChanged) isSignalSource() {}

type SignalSourcePushToTalkDelay struct {
	Source Source
	Delay  int64
}

func (SignalSourcePushToTalkDelay) isSignal()       {}
func (SignalSourcePushToTalkDelay) isSignalSource() {}

type SignalSourceEnable struct {
	Source  Source
	Enabled bool
}

func (SignalSourceEnable) isSignal()       {}
func (SignalSourceEnable) isSignalSource() {}

type SignalSourceRename struct {
	Source   Source
	NewName  string
	PrevName string
}

func (SignalSourceRename) isSignal()       {}
func (SignalSourceRename) isSignalSource() {}

type SignalSourceVolume struct {
	Source Source
	Volume float64
}

func (SignalSourceVolume) isSignal()       {}
func (SignalSourceVolume) isSignalSource() {}

type SignalSourceUpdateProperties struct {
	Source Source
}

func (SignalSourceUpdateProperties) isSignal()       {}
func (SignalSourceUpdateProperties) isSignalSource() {}

type SignalSourceUpdateFlags struct {
	Source Source
	Flags  int64
}

func (SignalSourceUpdateFlags) isSignal()       {}
func (SignalSourceUpdateFlags) isSignalSource() {}

type SignalSourceAudioSync struct {
	Source Source
	Offset int64
}

func (SignalSourceAudioSync) isSignal()       {}
func (SignalSourceAudioSync) isSignalSource() {}

type SignalSourceAudioBalance struct {
	Source  Source
	Balance float64
}

func (SignalSourceAudioBalance) isSignal()       {}
func (SignalSourceAudioBalance) isSignalSource() {}

type SignalSourceAudioMixers struct {
	Source Source
	Mixers int64
}

func (SignalSourceAudioMixers) isSignal()       {}
func (SignalSourceAudioMixers) isSignalSource() {}

type SignalSourceAudioActivate struct {
	Source Source
}

func (SignalSourceAudioActivate) isSignal()       {}
func (SignalSourceAudioActivate) isSignalSource() {}

type SignalSourceAudioDeactivate struct {
	Source Source
}

func (SignalSourceAudioDeactivate) isSignal()       {}
func (SignalSourceAudioDeactivate) isSignalSource() {}

type SignalSourceFilterAdd struct {
	Source Source
	Filter Filter
}

func (SignalSourceFilterAdd) isSignal()       {}
func (SignalSourceFilterAdd) isSignalSource() {}

type SignalSourceFilterRemove struct {
	Source Source
	Filter Filter
}

func (SignalSourceFilterRemove) isSignal()       {}
func (SignalSourceFilterRemove) isSignalSource() {}

type SignalSourceReorderFilters struct {
	Source Source
}

func (SignalSourceReorderFilters) isSignal()       {}
func (SignalSourceReorderFilters) isSignalSource() {}

type SignalSourceTransitionStart struct {
	Source Source
}

func (SignalSourceTransitionStart) isSignal()       {}
func (SignalSourceTransitionStart) isSignalSource() {}

type SignalSourceTransitionVideoStop struct {
	Source Source
}

func (SignalSourceTransitionVideoStop) isSignal()       {}
func (SignalSourceTransitionVideoStop) isSignalSource() {}

type SignalSourceTransitionStop struct {
	Source Source
}

func (SignalSourceTransitionStop) isSignal()       {}
func (SignalSourceTransitionStop) isSignalSource() {}

type SignalSourceMediaStarted struct {
	Source Source
}

func (SignalSourceMediaStarted) isSignal()       {}
func (SignalSourceMediaStarted) isSignalSource() {}

type SignalSourceMediaEnded struct {
	Source Source
}

func (SignalSourceMediaEnded) isSignal()       {}
func (SignalSourceMediaEnded) isSignalSource() {}

type SignalSourceMediaPause struct {
	Source Source
}

func (SignalSourceMediaPause) isSignal()       {}
func (SignalSourceMediaPause) isSignalSource() {}

type SignalSourceMediaPlay struct {
	Source Source
}

func (SignalSourceMediaPlay) isSignal()       {}
func (SignalSourceMediaPlay) isSignalSource() {}

type SignalSourceMediaRestart struct {
	Source Source
}

func (SignalSourceMediaRestart) isSignal()       {}
func (SignalSourceMediaRestart) isSignalSource() {}

type SignalSourceMediaStopped struct {
	Source Source
}

func (SignalSourceMediaStopped) isSignal()       {}
func (SignalSourceMediaStopped) isSignalSource() {}

type SignalSourceMediaNext struct {
	Source Source
}

func (SignalSourceMediaNext) isSignal()       {}
func (SignalSourceMediaNext) isSignalSource() {}

type SignalSourceMediaPrevious struct {
	Source Source
}

func (SignalSourceMediaPrevious) isSignal()       {}
func (SignalSourceMediaPrevious) isSignalSource() {}

type SignalImageSlideShowSlideChanged struct {
	Index int64
	Path  string
}

func (SignalImageSlideShowSlideChanged) isSignal()       {}
func (SignalImageSlideShowSlideChanged) isSignalSource() {}

type SignalSourceHoooked struct {
	Source     Source
	Title      string
	Class      string
	Executable string
}

func (SignalSourceHoooked) isSignal()       {}
func (SignalSourceHoooked) isSignalSource() {}

type SignalSourceUnhooked struct {
	Source Source
}

func (SignalSourceUnhooked) isSignal()       {}
func (SignalSourceUnhooked) isSignalSource() {}
