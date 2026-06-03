package alpm

// #include <stdint.h>
import "C"
import (
	"fmt"
	"math"
)

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
