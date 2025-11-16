package sources

import "github.com/dusk125/gobs"

type MediaSource struct {
	gobs.Source
}

func (ms MediaSource) Restart() {
	cd := gobs.CallData{}
	defer cd.Free()

	ms.ProcHandler().Call("restart", cd)
}

func (ms MediaSource) Duration() int64 {
	cd := gobs.CallData{}
	defer cd.Free()

	ms.ProcHandler().Call("get_duration", cd)
	return cd.Int("duration")
}

func (ms MediaSource) Frames() int64 {
	cd := gobs.CallData{}
	defer cd.Free()

	ms.ProcHandler().Call("get_nb_frames", cd)
	return cd.Int("num_frames")
}
