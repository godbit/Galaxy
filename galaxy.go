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
	tMax = 16 * time.Hour * 24
)

type Point struct {
	X float64
	Y float64
}

type Event struct {
	S Point
	T time.Time
}

func Cluster(events []Event) (Ns, N2s, Nt, N2t, X int) {
	// Input: D and T
	fmt.Println("\n\n====================================================")
	fmt.Println("Init space time cluster calculation")

	// Distance matches (first and second order)
	Ns = 0
	N2s = 0

	// Time matches (first and second order)
	Nt = 0
	N2t = 0

	// Both matching
	X = 0

	startTime := time.Now()

	for i := range events {
		if i%100 == 0 && i != 0 {
			fmt.Printf("%d features complete", i)
			fmt.Println("Time elapsed:", time.Since(startTime))
		}
		for j := range events {
			if i == j {
				// this is just to eliminate the self-pairing, not that we are still double counting, i.e., both ij and ji are counted which we later have to normalize
				continue
			}
			if dDiff(events[i].S, events[j].S) <= dMax {
				Ns++
			}
			if tDiff(events[i].T, events[j].T) <= tMax {
				Nt++
			}
			if dDiff(events[i].S, events[j].S) <= dMax && tDiff(events[i].T, events[j].T) <= tMax {
				X++
			}
			for k := range events {
				// the second order terms are also only double counted because the join of the pairs is only considered on j
				if i == k || j == k {
					continue
					// this is just to eliminate the self-pairing
				}
				if dDiff(events[i].S, events[j].S) <= dMax && dDiff(events[j].S, events[k].S) <= dMax {
					N2s++
				}
				if tDiff(events[i].T, events[j].T) <= tMax && tDiff(events[j].T, events[k].T) <= tMax {
					N2t++
				}
			}
		}
	}
	// normalize for double counting
	Ns = Ns / 2
	N2s = N2s / 2
	Nt = Nt / 2
	N2t = N2t / 2
	X = X / 2

	return Ns, N2s, Nt, N2t, X
}

func dDiff(a, b Point) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}

func tDiff(a, b time.Time) time.Duration {
	if a.Before(b) {
		return b.Sub(a)
	}
	return a.Sub(b)
}
