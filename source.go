package gobs

// #include <obs/obs.h>
import "C"

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

// Creates a source of the specified type with the specified settings.
// The “source” context is used for anything related to presenting or modifying video/audio.
// Use Source.Release() to release it.
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

func (s Source) Valid() bool {
	return s.c != nil
}

func (s Source) Release() {
	// #cgo noescape obs_source_release
	// #cgo nocallback obs_source_release
	C.obs_source_release(s.c)
}

func (s Source) Remove() {
	// #cgo noescape obs_source_remove
	// #cgo nocallback obs_source_remove
	C.obs_source_remove(s.c)
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

func (s Source) VideoRender() {
	// #cgo noescape obs_source_video_render
	// #cgo nocallback obs_source_video_render
	C.obs_source_video_render(s.c)
}

func (s Source) IncShowing() {
	// #cgo noescape obs_source_inc_showing
	// #cgo nocallback obs_source_inc_showing
	C.obs_source_inc_showing(s.c)
}

func (s Source) DecShowing() {
	// #cgo noescape obs_source_dec_showing
	// #cgo nocallback obs_source_dec_showing
	C.obs_source_dec_showing(s.c)
}

func (s Source) Width() uint32 {
	// #cgo noescape obs_source_get_width
	// #cgo nocallback obs_source_get_width
	return uint32(C.obs_source_get_width(s.c))
}

func (s Source) Height() uint32 {
	// #cgo noescape obs_source_get_height
	// #cgo nocallback obs_source_get_height
	return uint32(C.obs_source_get_height(s.c))
}

func (s Source) ProcHandler() ProcHandler {
	return ProcHandler{C.obs_source_get_proc_handler(s.c)}
}

func (s Source) SignalHandler() SignalHandler {
	return SignalHandler{C.obs_source_get_signal_handler(s.c)}
}

func (s Source) Volume() float32 {
	// #cgo noescape obs_source_get_volume
	// #cgo nocallback obs_source_get_volume
	return float32(C.obs_source_get_volume(s.c))
}

func (s Source) SetVolume(vol float32) {
	// #cgo noescape obs_source_set_volume
	// #cgo nocallback obs_source_set_volume
	C.obs_source_set_volume(s.c, C.float(vol))
}

func (s Source) Muted() bool {
	// #cgo noescape obs_source_muted
	// #cgo nocallback obs_source_muted
	return bool(C.obs_source_muted(s.c))
}

func (s Source) SetMuted(muted bool) {
	// #cgo noescape obs_source_set_muted
	// #cgo nocallback obs_source_set_muted
	C.obs_source_set_muted(s.c, C.bool(muted))
}
