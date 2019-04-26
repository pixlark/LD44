package main

type RGBA struct {
	R float32
	G float32
	B float32
	A float32
}

type HSV struct {
	H float32
	S float32
	V float32
}

func (this HSV) Rgba() RGBA {
	var hh, p, q, t, ff float32
	var i   int
	var out RGBA
	out.A = 1.0

	if this.S <= 0.0 {
		out.R = this.V
		out.G = this.V
		out.B = this.V
		return out
	}
	
	hh = this.H
	if (hh >= 360.0) {
		hh = 0.0
	}
	hh /= 60.0
	i = int(hh)
	ff = hh - float32(i)
	p = this.V * (1.0 - this.S)
	q = this.V * (1.0 - (this.S * ff))
	t = this.V * (1.0 - (this.S * (1.0 - ff)))

	switch i {
	case 0:
		out.R = this.V
		out.G = t
		out.B = p
	case 1:
		out.R = q
		out.G = this.V
		out.B = p
	case 2:
		out.R = p
		out.G = this.V
		out.B = t
	case 3:
		out.R = p
		out.G = q
		out.B = this.V
	case 4:
		out.R = t
		out.G = p
		out.B = this.V
	default:
		out.R = this.V
		out.G = p
		out.B = q
	}
	return out
}
