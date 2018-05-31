// Package galaxy calculates space-time correlations to the first and second
// order.
package galaxy

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/karlek/progress/barcli"
	"github.com/pkg/errors"
)

type Point struct {
	X float64
	Y float64
}

type Event struct {
	S Point
	T int64 //time.Time
}

func Cluster(ctx context.Context, events []Event, dMax float64, tMax int64, verbose bool) (Ns, N2s, Nt, N2t, X int) {
	// Input: D and T
	if verbose {
		fmt.Println("Init space time cluster calculation")
	}

	// Convert tMax from days to nanoseconds
	tMax = tMax*int64(time.Hour)*24 + int64(time.Hour)

	// Distance matches (first and second order)
	Ns = 0
	N2s = 0

	// Time matches (first and second order)
	Nt = 0
	N2t = 0

	// Both matching
	X = 0

	const nworkers = 4
	partSize := len(events) / nworkers
	c := make(chan Result)

	bar, err := barcli.New(len(events))
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	for i := 0; i < nworkers; i++ {
		imin := i * partSize
		imax := (i + 1) * partSize
		if imax >= len(events) {
			imax = len(events)
		}
		go inner(ctx, imin, imax, dMax, tMax, events, bar, verbose, c)
	}
	for i := 0; i < nworkers; i++ {
		result := <-c
		Ns += result.Ns
		N2s += result.N2s
		Nt += result.Nt
		N2t += result.N2t
		X += result.X
	}
	if verbose {
		// Print last status update.
		bar.Print()
	}

	// normalize for double counting
	Ns = Ns / 2
	N2s = N2s / 2
	Nt = Nt / 2
	N2t = N2t / 2
	X = X / 2

	return Ns, N2s, Nt, N2t, X
}

type Result struct {
	Ns  int
	N2s int
	Nt  int
	N2t int
	X   int
}

func inner(ctx context.Context, imin, imax int, dMax float64, tMax int64, events []Event, bar *barcli.Bar, verbose bool, c chan Result) {
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
			sdiff := dDiff(events[i].S, events[j].S)
			if sdiff <= dMax {
				result.Ns++
			}
			tdiff := tDiff(events[i].T, events[j].T)
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
				if sdiff <= dMax && dDiff(events[j].S, events[k].S) <= dMax {
					result.N2s++
				}
				if tdiff <= tMax && tDiff(events[j].T, events[k].T) <= tMax {
					result.N2t++
				}
			}
		}
	}

	c <- result
}

func dDiff(a, b Point) float64 {
	xdiff := a.X - b.X
	ydiff := a.Y - b.Y
	return math.Sqrt(xdiff*xdiff + ydiff*ydiff)
}

func tDiff(a, b int64) int64 {
	if a < b {
		return b - a
	}
	return a - b
}
