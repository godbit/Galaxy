//+build cgo

package galaxy

import (
	"context"
	"fmt"
	"log"
	"unsafe"

	"github.com/karlek/progress/barcli"
)

// #cgo LDFLAGS: -lm
//
// #include <stdint.h>
// #include <math.h>
//
// typedef struct {
//    double X;
//    double Y;
// } Point;
//
// double ddiff_c(Point a, Point b) {
//    double xdiff = a.X - b.X;
//    double ydiff = a.Y - b.Y;
//    return sqrt(xdiff*xdiff + ydiff*ydiff);
// }
//
// int64_t tdiff_c(int64_t a, int64_t b) {
//    if (a < b) {
//       return b - a;
//    }
//    return a - b;
// }
import "C"

func init() {
	log.Println("inner loop in C")
}

func inner(ctx context.Context, imin, imax int, events []Event, bar *barcli.Bar, verbose bool, c chan Result) {
	var result Result

	for i := imin; i < imax; i++ {
		// Send partial results on interrupt.
		select {
		case <-ctx.Done():
			fmt.Println()
			log.Printf("sending partial results for i = %d (%d iterations) in range [%d, %d)", i, i-imin, imin, imax)
			c <- result
			return
		default:
		}

		if verbose {
			bar.Inc()
			// Only print status updates from the first worker Go routine.
			if imin == 0 {
				bar.Print()
			}
		}

		for j := range events {
			if i == j {
				// this is just to eliminate the self-pairing, not that we are still double counting, i.e., both ij and ji are counted which we later have to normalize
				continue
			}
			sdiff := float64(C.ddiff_c(*(*C.Point)(unsafe.Pointer(&events[i].S)), *(*C.Point)(unsafe.Pointer(&events[j].S))))
			if sdiff <= dMax {
				result.Ns++
			}
			tdiff := int64(C.tdiff_c(C.int64_t(events[i].T), C.int64_t(events[j].T)))
			if tdiff <= tMax {
				result.Nt++
			}
			if sdiff <= dMax && tdiff <= tMax {
				result.X++
			}
			for k := range events {
				// the second order terms are also only double counted because the join of the pairs is only considered on j
				if i == k || j == k {
					continue
					// this is just to eliminate the self-pairing
				}
				if sdiff <= dMax && float64(C.ddiff_c(*(*C.Point)(unsafe.Pointer(&events[j].S)), *(*C.Point)(unsafe.Pointer(&events[k].S)))) <= dMax {
					result.N2s++
				}
				if tdiff <= tMax && int64(C.tdiff_c(C.int64_t(events[j].T), C.int64_t(events[k].T))) <= tMax {
					result.N2t++
				}
			}
		}
	}

	c <- result
}
