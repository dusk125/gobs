package gobs

// #include <obs/obs.h>
import "C"
import "unsafe"

type Source struct {
	c *C.obs_source_t
}

func SourceDefaults(id string) Data {
	cid := fromString(id).cptr()

	data := obs_data{C.obs_get_source_defaults(cid)}
	defer data.Release()
	return data.MapWithDefaults()
}

// Creates a source of the specified type with the specified settings.
// The “source” context is used for anything related to presenting or modifying video/audio.
// Use Source.Release() to release it.
func SourceCreate(id string, name string, settings Data, hotkey Data) (Source, Releaser) {
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

	source := Source{C.obs_source_create(fromString(id).cptr(), fromString(name).cptr(), obsSettings.obs_data_t, obsHotkey.obs_data_t)}
	if source.c == nil {
		panic("failed to create source")
	}
	return source, source.Release
}

func (s Source) Valid() bool {
	return s.c != nil
}

func (s Source) Release() {
	C.obs_source_release(s.c)
}

func (s Source) Remove() {
	C.obs_source_remove(s.c)
}

func (s Source) Update(settings Data) {
	d := settings.obs_data()
	defer d.Release()

	C.obs_source_update(s.c, d.obs_data_t)
}

func (s Source) Name() string {
	return C.GoString(C.obs_source_get_name(s.c))
}

func (s Source) VideoRender() {
	C.obs_source_video_render(s.c)
}

func (s Source) IncShowing() {
	C.obs_source_inc_showing(s.c)
}

func (s Source) DecShowing() {
	C.obs_source_dec_showing(s.c)
}

func (s Source) Width() uint32 {
	return uint32(C.obs_source_get_width(s.c))
}

func (s Source) Height() uint32 {
	return uint32(C.obs_source_get_height(s.c))
}

func (s Source) ProcHandler() ProcHandler {
	return ProcHandler{C.obs_source_get_proc_handler(s.c)}
}

func (s Source) SignalHandler() SignalHandler {
	return SignalHandler{C.obs_source_get_signal_handler(s.c)}
}

func (s Source) Volume() float32 {
	return float32(C.obs_source_get_volume(s.c))
}

func (s Source) SetVolume(vol float32) {
	C.obs_source_set_volume(s.c, C.float(vol))
}

func (s Source) Muted() bool {
	return bool(C.obs_source_muted(s.c))
}

func (s Source) SetMuted(muted bool) {
	C.obs_source_set_muted(s.c, C.bool(muted))
}

func (s Source) PressButton(name string) bool {
	props := C.obs_source_properties(s.c)
	if props == nil {
		return false
	}
	defer C.obs_properties_destroy(props)

	p := C.obs_properties_get(props, fromString(name).cptr())
	if p == nil {
		return false
	}
	C.obs_property_button_clicked(p, unsafe.Pointer(s.c))
	return true
}
