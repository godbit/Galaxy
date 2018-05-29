// Package galaxy calculates space-time correlations to the first and second
// order.
package galaxy

import (
	"fmt"
	"math"
	"time"
)

const (
	// max distance in meters
	dMax = 1800
	// max temporal difference in days
	tMax = int64(16*time.Hour*24 + time.Hour)
)

type Point struct {
	X float64
	Y float64
}

type Event struct {
	S Point
	T int64 //time.Time
}

// debug output
const dbg = false

func Cluster(events []Event) (Ns, N2s, Nt, N2t, X int) {
	// Input: D and T
	if dbg {
		fmt.Println("\n\n====================================================")
		fmt.Println("Init space time cluster calculation")
	}

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
	for i := 0; i < nworkers; i++ {
		imin := i * partSize
		imax := (i + 1) * partSize
		if imax >= len(events) {
			imax = len(events)
		}
		go inner(imin, imax, events, c)
	}
	for i := 0; i < nworkers; i++ {
		result := <-c
		Ns += result.Ns
		N2s += result.N2s
		Nt += result.Nt
		N2t += result.N2t
		X += result.X
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

func inner(imin, imax int, events []Event, c chan Result) {
	var result Result

	startTime := time.Now()

	for i := imin; i < imax; i++ {
		if dbg {
			if i%100 == 0 && i != 0 {
				fmt.Printf("%d features complete", i)
				fmt.Println("Time elapsed:", time.Since(startTime))
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
