package alpm

// #include <limits.h>
// #include <stdint.h>
import "C"
import (
	"fmt"
	"math"
)

func c2goInt(x C.int) int {
	if C.INT_MIN <= math.MinInt {
		gmin := math.MinInt
		cmin := C.int(gmin)
		if x <= cmin {
			panic(fmt.Sprintf("C.int value %d overflows int", x))
		}
	}
	if C.INT_MAX >= math.MaxInt {
		gmax := math.MaxInt
		cmax := C.int(gmax)
		if x >= cmax {
			panic(fmt.Sprintf("C.int value %d overflows int", x))
		}
	}
	return int(x)
}

func c2goSize(x C.size_t) int {
	if C.SIZE_MAX >= math.MaxInt {
		gmax := math.MaxInt
		cmax := C.size_t(gmax)
		if x >= cmax {
			panic(fmt.Sprintf("C.size_t value %d overflows int", x))
		}
	}
	return int(x)
}
