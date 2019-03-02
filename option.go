package stlminus

// Option is for functional options.
type Option func(*stl)

// OuterLoop sets outer var.
func OuterLoop(outer int) Option {
	return func(args *stl) {
		args.outer = outer
		if outer < 0 {
			args.outer = outer * -1
		}
	}
}

// InnerLoop sets inner var.
func InnerLoop(inner int) Option {
	return func(args *stl) {
		args.inner = inner
		if inner < 0 {
			args.inner = inner * -1
		}
	}
}

// SWindow sets sWindow var.
func SWindow(sWindow int) Option {
	return func(args *stl) {
		args.sWindow = sWindow
		if sWindow < 0 {
			args.sWindow = sWindow * -1
		}
	}
}

// TWindow sets tWindow var.
func TWindow(tWindow int) Option {
	return func(args *stl) {
		args.tWindow = tWindow
		if tWindow < 0 {
			args.tWindow = tWindow * -1
		}
	}
}

// LWindow sets lWindow var.
func LWindow(lWindow int) Option {
	return func(args *stl) {
		args.lWindow = lWindow
		if lWindow < 0 {
			args.lWindow = lWindow * -1
		}
	}
}

// SDegree sets sDegree var.
func SDegree(sDegree int) Option {
	return func(args *stl) {
		args.sDegree = sDegree
		if sDegree < 0 || sDegree > 2 {
			args.sDegree = 1
		}
	}
}

// TDegree sets tDegree var.
func TDegree(tDegree int) Option {
	return func(args *stl) {
		args.tDegree = tDegree
		if tDegree < 0 || tDegree > 2 {
			args.tDegree = 1
		}
	}
}

// LDegree sets lDegree var.
func LDegree(lDegree int) Option {
	return func(args *stl) {
		args.lDegree = lDegree
		if lDegree < 0 || lDegree > 2 {
			args.lDegree = 1
		}
	}
}

// SJump sets sJump var.
func SJump(sJump int) Option {
	return func(args *stl) {
		args.sJump = sJump
		if sJump < 0 {
			args.sJump = sJump * -1
		}
	}
}

// TJump sets tJump var.
func TJump(tJump int) Option {
	return func(args *stl) {
		args.tJump = tJump
		if tJump < 0 {
			args.tJump = tJump * -1
		}
	}
}

// LJump sets lJump var.
func LJump(lJump int) Option {
	return func(args *stl) {
		args.lJump = lJump
		if lJump < 0 {
			args.lJump = lJump * -1
		}
	}
}

// CritFreq sets critFreq var.
func CritFreq(critFreq float64) Option {
	return func(args *stl) {
		args.critFreq = critFreq
		if critFreq < 0 {
			args.critFreq = critFreq * -1
		}
	}
}
