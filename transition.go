package gobs

// #include <obs/obs.h>
import "C"
import "time"

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

	data := obs_data{C.obs_get_source_defaults(cid)}
	defer data.Release()
	return data.MapWithDefaults()
}

func TransitionCreate(id, name string, settings, hotkeys Data) (Transition, Releaser) {
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

	t := Transition{Source{C.obs_source_create(fromString(id).cptr(), fromString(name).cptr(), obsSettings.obs_data_t, obsHotkeys.obs_data_t)}}
	if t.c == nil {
		panic("failed to create transition")
	}
	return t, t.Release
}

func (t Transition) Start(mode TransitionMode, duration time.Duration, dest Source) bool {
	return bool(C.obs_transition_start(t.c, C.enum_obs_transition_mode(mode), C.uint32_t(duration.Milliseconds()), dest.c))
}

func (t Transition) Set(source Source) {
	C.obs_transition_set(t.c, source.c)
}

func (t Transition) Clear() {
	C.obs_transition_clear(t.c)
}

func (t Transition) GetActiveSource() (Source, Releaser) {
	src := Source{C.obs_transition_get_active_source(t.c)}
	return src, src.Release
}

func (t Transition) GetSource(target TransitionTarget) (Source, Releaser) {
	src := Source{C.obs_transition_get_source(t.c, C.enum_obs_transition_target(target))}
	return src, src.Release
}

func (t Transition) SetScaleType(st TransitionScaleType) {
	C.obs_transition_set_scale_type(t.c, C.enum_obs_transition_scale_type(st))
}

func (t Transition) ScaleType() TransitionScaleType {
	return TransitionScaleType(C.obs_transition_get_scale_type(t.c))
}

func (t Transition) SetAlignment(alignment uint32) {
	C.obs_transition_set_alignment(t.c, C.uint32_t(alignment))
}

func (t Transition) Alignment() uint32 {
	return uint32(C.obs_transition_get_alignment(t.c))
}

func (t Transition) SetSize(cx, cy uint32) {
	C.obs_transition_set_size(t.c, C.uint32_t(cx), C.uint32_t(cy))
}

func (t Transition) Size() (uint32, uint32) {
	var cx, cy C.uint32_t
	C.obs_transition_get_size(t.c, &cx, &cy)
	return uint32(cx), uint32(cy)
}

func (t Transition) EnableFixed(enable bool, duration time.Duration) {
	C.obs_transition_enable_fixed(t.c, C.bool(enable), C.uint32_t(duration.Milliseconds()))
}

func (t Transition) Fixed() bool {
	return bool(C.obs_transition_fixed(t.c))
}

func (t Transition) SetManualTime(tm float32) {
	C.obs_transition_set_manual_time(t.c, C.float(tm))
}

func (t Transition) SetManualTorque(torque, clamp float32) {
	C.obs_transition_set_manual_torque(t.c, C.float(torque), C.float(clamp))
}

func (t Transition) GetTime() float32 {
	return float32(C.obs_transition_get_time(t.c))
}

func (t Transition) ForceStop() {
	C.obs_transition_force_stop(t.c)
}

func (t Transition) IsActive() bool {
	return bool(C.obs_transition_is_active(t.c))
}
