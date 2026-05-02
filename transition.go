package gobs

// #include <obs/obs.h>
import "C"

type Transition struct {
	Source
}

type TransitionTarget uint32

const (
	TRANSITION_SOURCE_A TransitionTarget = C.OBS_TRANSITION_SOURCE_A
	TRANSITION_SOURCE_B TransitionTarget = C.OBS_TRANSITION_SOURCE_B
)

type TransitionMode uint32

const (
	TRANSITION_MODE_AUTO   TransitionMode = C.OBS_TRANSITION_MODE_AUTO
	TRANSITION_MODE_MANUAL TransitionMode = C.OBS_TRANSITION_MODE_MANUAL
)

type TransitionScaleType uint32

const (
	TRANSITION_SCALE_MAX_ONLY TransitionScaleType = C.OBS_TRANSITION_SCALE_MAX_ONLY
	TRANSITION_SCALE_ASPECT   TransitionScaleType = C.OBS_TRANSITION_SCALE_ASPECT
	TRANSITION_SCALE_STRETCH  TransitionScaleType = C.OBS_TRANSITION_SCALE_STRETCH
)

func TransitionDefaults(id string) Data {
	cid := fromString(id).cptr()

	// #cgo noescape obs_get_source_defaults
	// #cgo nocallback obs_get_source_defaults
	data := obs_data{C.obs_get_source_defaults(cid)}
	defer data.Release()
	return data.MapWithDefaults()
}

func TransitionCreate(id, name string, settings, hotkeys Data) Transition {
	var (
		obsSettings obs_data
		obsHotkeys  obs_data
	)

	if settings != nil {
		obsSettings = settings.obs_data()
		defer obsSettings.Release()
	}
	if hotkeys != nil {
		obsHotkeys = hotkeys.obs_data()
		defer obsHotkeys.Release()
	}

	// #cgo noescape obs_source_create
	// #cgo nocallback obs_source_create
	t := Transition{Source{C.obs_source_create(fromString(id).cptr(), fromString(name).cptr(), obsSettings.obs_data_t, obsHotkeys.obs_data_t)}}
	if t.c == nil {
		panic("failed to create transition")
	}
	return t
}

func (t Transition) Start(mode TransitionMode, durationMs uint32, dest Source) bool {
	// #cgo noescape obs_transition_start
	// #cgo nocallback obs_transition_start
	return bool(C.obs_transition_start(t.c, C.enum_obs_transition_mode(mode), C.uint32_t(durationMs), dest.c))
}

func (t Transition) Set(source Source) {
	// #cgo noescape obs_transition_set
	// #cgo nocallback obs_transition_set
	C.obs_transition_set(t.c, source.c)
}

func (t Transition) Clear() {
	// #cgo noescape obs_transition_clear
	// #cgo nocallback obs_transition_clear
	C.obs_transition_clear(t.c)
}

func (t Transition) GetActiveSource() Source {
	// #cgo noescape obs_transition_get_active_source
	// #cgo nocallback obs_transition_get_active_source
	return Source{C.obs_transition_get_active_source(t.c)}
}

func (t Transition) GetSource(target TransitionTarget) Source {
	// #cgo noescape obs_transition_get_source
	// #cgo nocallback obs_transition_get_source
	return Source{C.obs_transition_get_source(t.c, C.enum_obs_transition_target(target))}
}

func (t Transition) SetScaleType(st TransitionScaleType) {
	// #cgo noescape obs_transition_set_scale_type
	// #cgo nocallback obs_transition_set_scale_type
	C.obs_transition_set_scale_type(t.c, C.enum_obs_transition_scale_type(st))
}

func (t Transition) ScaleType() TransitionScaleType {
	// #cgo noescape obs_transition_get_scale_type
	// #cgo nocallback obs_transition_get_scale_type
	return TransitionScaleType(C.obs_transition_get_scale_type(t.c))
}

func (t Transition) SetAlignment(alignment uint32) {
	// #cgo noescape obs_transition_set_alignment
	// #cgo nocallback obs_transition_set_alignment
	C.obs_transition_set_alignment(t.c, C.uint32_t(alignment))
}

func (t Transition) Alignment() uint32 {
	// #cgo noescape obs_transition_get_alignment
	// #cgo nocallback obs_transition_get_alignment
	return uint32(C.obs_transition_get_alignment(t.c))
}

func (t Transition) SetSize(cx, cy uint32) {
	// #cgo noescape obs_transition_set_size
	// #cgo nocallback obs_transition_set_size
	C.obs_transition_set_size(t.c, C.uint32_t(cx), C.uint32_t(cy))
}

func (t Transition) Size() (uint32, uint32) {
	var cx, cy C.uint32_t
	// #cgo noescape obs_transition_get_size
	// #cgo nocallback obs_transition_get_size
	C.obs_transition_get_size(t.c, &cx, &cy)
	return uint32(cx), uint32(cy)
}

func (t Transition) EnableFixed(enable bool, durationMs uint32) {
	// #cgo noescape obs_transition_enable_fixed
	// #cgo nocallback obs_transition_enable_fixed
	C.obs_transition_enable_fixed(t.c, C.bool(enable), C.uint32_t(durationMs))
}

func (t Transition) Fixed() bool {
	// #cgo noescape obs_transition_fixed
	// #cgo nocallback obs_transition_fixed
	return bool(C.obs_transition_fixed(t.c))
}

func (t Transition) SetManualTime(tm float32) {
	// #cgo noescape obs_transition_set_manual_time
	// #cgo nocallback obs_transition_set_manual_time
	C.obs_transition_set_manual_time(t.c, C.float(tm))
}

func (t Transition) SetManualTorque(torque, clamp float32) {
	// #cgo noescape obs_transition_set_manual_torque
	// #cgo nocallback obs_transition_set_manual_torque
	C.obs_transition_set_manual_torque(t.c, C.float(torque), C.float(clamp))
}

func (t Transition) GetTime() float32 {
	// #cgo noescape obs_transition_get_time
	// #cgo nocallback obs_transition_get_time
	return float32(C.obs_transition_get_time(t.c))
}

func (t Transition) ForceStop() {
	// #cgo noescape obs_transition_force_stop
	// #cgo nocallback obs_transition_force_stop
	C.obs_transition_force_stop(t.c)
}
