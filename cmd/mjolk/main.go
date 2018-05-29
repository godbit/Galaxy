package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	galaxy "github.com/godbit/Galaxy"
	"github.com/godbit/Galaxy/knox"
)

func main() {
	flag.Parse()
	for _, jsonPath := range flag.Args() {
		events, err := galaxy.ParseFile(jsonPath)
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

		std := math.Sqrt(V)

		fmt.Println("\nStatistics:")
		fmt.Println("N:", N)
		fmt.Println("E:", E)
		fmt.Println("V:", V)
		fmt.Println("Std:", std)
		fmt.Println("\nZ-score:")
		diff := math.Abs(float64(X) - float64(E))
		Z := diff / std
		fmt.Println("Z:", Z)
	}
}
