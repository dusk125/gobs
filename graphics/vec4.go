package graphics

type Vec4 [4]float32

func (v Vec4) X() float32 {
	return v[0]
}

func (v Vec4) Y() float32 {
	return v[1]
}

func (v Vec4) Z() float32 {
	return v[2]
}

func (v Vec4) W() float32 {
	return v[3]
}
