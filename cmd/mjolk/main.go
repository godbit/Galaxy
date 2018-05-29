package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/godbit/Galaxy/knox"
	"github.com/pkg/errors"
	"github.com/simplereach/timeutils"
)

func main() {
	flag.Parse()
	for _, jsonPath := range flag.Args() {
		events, err := parseJson(jsonPath)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		Ns, N2s, Nt, N2t, X := calcSpaceTimeCluster(events)

		fmt.Println("\nCounts:")
		fmt.Println("Ns: ", Ns)
		fmt.Println("N2s: ", N2s)
		fmt.Println("Nt: ", Nt)
		fmt.Println("N2t: ", N2t)
		fmt.Println("X: ", X)
		n := len(events)
		fmt.Println("n: ", n)

		N, E, V := knox.Test(float64(Ns), float64(N2s), float64(Nt), float64(N2t), float64(X), float64(n))
		fmt.Println("N:", N)
		fmt.Println("E:", E)
		fmt.Println("V:", V)
	}
}

func parseEvent(data []interface{}) (Event, error) {
	coords := data[2].([]interface{})
	date, err := timeutils.ParseDateString(data[1].(string))
	if err != nil {
		return Event{}, errors.WithStack(err)
	}
	return Event{
		T: date,
		S: Point{X: coords[0].(float64), Y: coords[1].(float64)},
	}, nil
}

func parseJson(jsonPath string) ([]Event, error) {
	buf, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var data [][]interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, errors.WithStack(err)
	}
	var events []Event
	for _, d := range data {
		event, err := parseEvent(d)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		events = append(events, event)
	}
	return events, nil
}

type Point struct {
	X float64
	Y float64
}

type Event struct {
	S Point
	T time.Time
}

const (
	// max distance in meters
	dMax = 1800
	// max temporal difference in days
	tMax = 16 * time.Hour * 24
)

func calcSpaceTimeCluster(events []Event) (Ns, N2s, Nt, N2t, X int) {
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

	//N := len(events) * len(events)

	//E[X] = Nt * Ns / N
	//V[X] = foo // …the formula on page 545 in Kulldorff et al.
	// V[X] should be equal to E[X] for Poisson distribution and should be close to each other based on the calculations if the variable X is indeed approximately Poisson distributed.
	//pretty.Println(V[X])
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
