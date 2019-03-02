package stlminus

/*
extern void stl_(
    float [], 	// y
    int *,    	// n
    int *,    	// np
    int *,    	// ns
    int *,    	// nt
    int *,   	// nl
    int *,   	// isdeg
    int *,    	// itdeg
    int *,   	// ildeg
    int *,   	// nsjump
    int *,    	// ntjump
    int *,    	// nljump
    int *,    	// ni
    int *,    	// no
    float [],   // rw
    float [],   // season
    float [],   // trend
    float [][7] // work
);
*/
import "C"
import (
	"errors"
	"math"
)

type stl struct {
	outer    int
	inner    int
	sWindow  int
	lWindow  int
	tWindow  int
	sDegree  int
	lDegree  int
	tDegree  int
	sJump    int
	lJump    int
	tJump    int
	critFreq float64
}

// Decompose performs a STL decomposition.
func Decompose(series []float64, seasonality int, opts ...Option) (trend, seasonal, remainder []float64, err error) {
	sl := len(series)
	if sl < 11 {
		return nil, nil, nil, errors.New("series length must be at least 11")
	}
	if seasonality < 5 {
		return nil, nil, nil, errors.New("seasonality must be at least 5")
	}
	if sl <= 2*seasonality {
		return nil, nil, nil, errors.New("series length must be > 2 * seasonality")
	}

	// Default Options
	s := &stl{
		sWindow:  -1,
		tWindow:  -1,
		lWindow:  -1,
		sDegree:  1,
		tDegree:  1,
		lDegree:  1,
		sJump:    -1,
		tJump:    -1,
		lJump:    -1,
		inner:    2,
		outer:    1,
		critFreq: 0.05,
	}

	for _, opt := range opts {
		opt(s)
	}

	removeNaNByAvg(series)
	var series32 = make([]float32, sl)
	for k := range series {
		series32[k] = float32(series[k])
	}
	return s.decompose(series32, seasonality)
}

func removeNaNByAvg(arr []float64) {
	n := len(arr)
	for i := 0; i < n; i++ {
		if math.IsNaN(arr[i]) {
			sum := 0.0
			count := 0
			if i-1 >= 0 && !math.IsNaN(arr[i-1]) {
				sum += arr[i-1]
				count++
			}
			if i+1 < n && !math.IsNaN(arr[i+1]) {
				sum += arr[i+1]
				count++
			}
			if count != 0 {
				arr[i] = sum / float64(count)
			} else {
				arr[i] = 0.0
			}
		}
	}
}

func nextOdd(x float64) int {
	xx := int(math.Round(x))
	if xx%2 == 0 {
		return xx + 1
	}
	return xx
}

func (s *stl) decompose(series []float32, seasonality int) ([]float64, []float64, []float64, error) {
	l := len(series)
	var y = make([]float32, l)
	copy(y, series)
	var rw = make([]float32, l)
	var season = make([]float32, l)
	var trend = make([]float32, l)

	s.sWindow = 10*l + 1
	s.sDegree = 0
	s.sJump = int(math.Ceil(float64(s.sWindow) / 10.0))

	if s.tWindow == -1 {
		//s.tWindow = nextOdd(math.Ceil(1.5*float64(seasonality)/(1.0-1.5/float64(s.sWindow)) + 0.5))
		s.tWindow = calcTWindow(s.tDegree, s.sDegree, s.sWindow, seasonality, s.critFreq)
	} else {
		s.tWindow = nextOdd(float64(s.tWindow))
	}

	if s.lWindow == -1 {
		s.lWindow = nextOdd(float64(seasonality))
	} else {
		s.lWindow = nextOdd(float64(s.lWindow))
	}

	if s.tWindow > l {
		s.tWindow = nextOdd(float64(l - 2))
	}
	if s.lWindow > l {
		s.lWindow = nextOdd(float64(l - 2))
	}

	if s.sJump == -1 || s.sJump == 0 {
		s.sJump = int(math.Ceil(float64(s.sWindow) / 10.0))
	}
	if s.tJump == -1 || s.tJump == 0 {
		s.tJump = int(math.Ceil(float64(s.tWindow) / 10.0))
	}
	if s.lJump == -1 || s.lJump == 0 {
		s.lJump = int(math.Ceil(float64(s.lWindow) / 10.0))
	}

	n := C.int(len(series))
	np := C.int(seasonality)
	var work = make([][7]C.float, n+2*np)
	for i := 0; i < l+2*seasonality; i++ {
		for j := 0; j < 7; j++ {
			work[i][j] = 0.0
		}
	}

	// nt = smallest odd integer greater than or equal to (1.5*np) / (1-(1.5/ns)).
	// nl = smallest odd integer greater than or equal to np.
	// nsjump = ns/10;
	// ntjump = nt/10;
	// nljump = nl/10.
	// ni = 2 if robust=.false.; else ni=1.
	// no = 0 if robust=.false.; else robustness iterations are  carried
	// out until convergence of both seasonal and trend components, with
	// 15 iterations maximum.  Convergence occurs if the maximum changes
	// in  individual  seasonal  and  trend fits are less than 1% of the
	// component's range after the previous iteration.

	ns := C.int(s.sWindow)
	nt := C.int(s.tWindow)
	nl := C.int(s.lWindow)
	isdeg := C.int(s.sDegree)
	itdeg := C.int(s.tDegree)
	ildeg := C.int(s.lDegree)
	nsjump := C.int(s.sJump)
	ntjump := C.int(s.lJump)
	nljump := C.int(s.tJump)
	ni := C.int(s.inner)
	no := C.int(s.outer)

	C.stl_(
		(*C.float)(&y[0]),
		&n,
		&np,
		&ns,
		&nt,
		&nl,
		&isdeg,
		&itdeg,
		&ildeg,
		&nsjump,
		&ntjump,
		&nljump,
		&ni,
		&no,
		(*C.float)(&rw[0]),
		(*C.float)(&season[0]),
		(*C.float)(&trend[0]),
		&work[0],
	)

	var t64 = make([]float64, l)
	var s64 = make([]float64, l)
	var r64 = make([]float64, l)

	for i := 0; i < len(season); i++ {
		t64[i] = float64(trend[i])
		s64[i] = float64(season[i])
		r64[i] = float64(rw[i])
	}

	return t64, s64, r64, nil
}
