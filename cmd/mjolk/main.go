package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	galaxy "github.com/godbit/Galaxy"
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
		Ns, N2s, Nt, N2t, X := galaxy.Cluster(events)

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

func parseJson(jsonPath string) ([]galaxy.Event, error) {
	buf, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var data [][]interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, errors.WithStack(err)
	}
	var events []galaxy.Event
	for _, d := range data {
		event, err := parseEvent(d)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		events = append(events, event)
	}
	return events, nil
}

func parseEvent(data []interface{}) (galaxy.Event, error) {
	coords := data[2].([]interface{})
	date, err := timeutils.ParseDateString(data[1].(string))
	if err != nil {
		return galaxy.Event{}, errors.WithStack(err)
	}
	return galaxy.Event{
		T: date,
		S: galaxy.Point{X: coords[0].(float64), Y: coords[1].(float64)},
	}, nil
}
