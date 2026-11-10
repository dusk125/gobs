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

func (src Source) Valid() bool {
	return src.c != nil
}

func (src Source) Release() {
	C.obs_source_release(src.c)
}

func (src Source) Remove() {
	C.obs_source_remove(src.c)
}

func (src Source) Update(settings Data) {
	d := settings.obs_data()
	defer d.Release()

	C.obs_source_update(src.c, d.obs_data_t)
}

func (src Source) Name() string {
	return C.GoString(C.obs_source_get_name(src.c))
}

func (src Source) VideoRender() {
	C.obs_source_video_render(src.c)
}

func (src Source) IncShowing() {
	C.obs_source_inc_showing(src.c)
}

func (src Source) DecShowing() {
	C.obs_source_dec_showing(src.c)
}

func (src Source) Width() uint32 {
	return uint32(C.obs_source_get_width(src.c))
}

func (src Source) Height() uint32 {
	return uint32(C.obs_source_get_height(src.c))
}

func (src Source) ProcHandler() ProcHandler {
	return ProcHandler{C.obs_source_get_proc_handler(src.c)}
}

func (src Source) SignalHandler() SignalHandler {
	return SignalHandler{C.obs_source_get_signal_handler(src.c)}
}

func (src Source) Volume() float32 {
	return float32(C.obs_source_get_volume(src.c))
}

func (src Source) SetVolume(vol float32) {
	C.obs_source_set_volume(src.c, C.float(vol))
}

func (src Source) Muted() bool {
	return bool(C.obs_source_muted(src.c))
}

func (src Source) SetMuted(muted bool) {
	C.obs_source_set_muted(src.c, C.bool(muted))
}

func (src Source) PressButton(name string) bool {
	props := C.obs_source_properties(src.c)
	if props == nil {
		return false
	}
	defer C.obs_properties_destroy(props)

	p := C.obs_properties_get(props, fromString(name).cptr())
	if p == nil {
		return false
	}
	C.obs_property_button_clicked(p, unsafe.Pointer(src.c))
	return true
}

func (s Source) GetRef() (Source, Releaser) {
	src := Source{C.obs_source_get_ref(s.c)}
	return src, src.Release
}
