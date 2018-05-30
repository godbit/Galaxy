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

func Cluster(ctx context.Context, events []Event, verbose bool) (Ns, N2s, Nt, N2t, X int) {
	// Input: D and T
	if verbose {
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
		go inner(ctx, imin, imax, events, bar, verbose, c)
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
