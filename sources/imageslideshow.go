package sources

import "github.com/dusk125/gobs"

type ImageSlideShow struct {
	gobs.Source
}

func (iss ImageSlideShow) CurrentIndex() int64 {
	cd := gobs.CallData{}
	defer cd.Free()

	iss.ProcHandler().Call("current_index", cd)

	return cd.Int("current_index")
}

func (iss ImageSlideShow) TotalFiles() int64 {
	cd := gobs.CallData{}
	defer cd.Free()

	iss.ProcHandler().Call("total_files", cd)

	return cd.Int("total_files")
}
