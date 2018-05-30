//+build cgo

package galaxy

import (
	"context"
	"log"
	"unsafe"

	"github.com/karlek/progress/barcli"
)

/*
#cgo CFLAGS: -O3
#cgo LDFLAGS: -lm

#include <stdint.h>
#include <math.h>

typedef struct {
	double X;
	double Y;
} Point;

double ddiff_c(Point a, Point b) {
	double xdiff = a.X - b.X;
	double ydiff = a.Y - b.Y;
	return sqrt(xdiff*xdiff + ydiff*ydiff);
}

int64_t tdiff_c(int64_t a, int64_t b) {
	if (a < b) {
		return b - a;
	}
	return a - b;
}

typedef struct {
	int Ns;
	int N2s;
	int Nt;
	int N2t;
	int X;
} Result;

typedef struct {
	Point S;
	int64_t T;
} Event;

#define Hour 3600000000000 // nanoseconds

#define D_MAX 1800
#define T_MAX (16*Hour*24 + Hour)

Result inner(int imin, int imax, Event *events, int len) {
	Result result;

	for (int i = imin; i < imax; i++) {
		for (int j = 0; j < len; j++) {
			if (i == j) {
				// this is just to eliminate the self-pairing, not that we are still double counting, i.e., both ij and ji are counted which we later have to normalize
				continue;
			}
			double sdiff = ddiff_c(events[i].S, events[j].S);
			if (sdiff <= D_MAX) {
				result.Ns++;
			}
			double tdiff = tdiff_c(events[i].T, events[j].T);
			if (tdiff <= T_MAX) {
				result.Nt++;
			}
			if (sdiff <= D_MAX && tdiff <= T_MAX) {
				result.X++;
			}
			for (int k = 0; k < len; k++) {
				// the second order terms are also only double counted because the join of the pairs is only considered on j
				if (i == k || j == k) {
					continue;
					// this is just to eliminate the self-pairing
				}
				if (sdiff <= D_MAX && ddiff_c(events[j].S, events[k].S) <= D_MAX) {
					result.N2s++;
				}
				if (tdiff <= T_MAX && tdiff_c(events[j].T, events[k].T) <= T_MAX) {
					result.N2t++;
				}
			}
		}
	}

	return result;
}
*/
import "C"

func init() {
	log.Println("inner loop in C")
}

func inner(ctx context.Context, imin, imax int, events []Event, bar *barcli.Bar, verbose bool, c chan Result) {
	r := C.inner(C.int(imin), C.int(imax), (*C.Event)(unsafe.Pointer(&events[0])), C.int(len(events)))
	result := Result{
		Ns:  int(r.Ns),
		N2s: int(r.N2s),
		Nt:  int(r.Nt),
		N2t: int(r.N2t),
		X:   int(r.X),
	}
	c <- result
}
