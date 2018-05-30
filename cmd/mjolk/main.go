package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/godbit/Galaxy/galaxy"
	"github.com/godbit/Galaxy/knox"
)

var cgo = false

func main() {
	var (
		verbose bool
	)
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()
	for _, jsonPath := range flag.Args() {
		events, err := galaxy.ParseFile(jsonPath)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		// Create context to interrupt calculation and still receive partial
		// results.
		ctx := context.TODO()
		if !cgo {
			var cancel func()
			ctx, cancel = context.WithCancel(context.Background())
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
			go func() {
				select {
				case <-sig:
					cancel()
					return
				}
			}()
		}

		Ns, N2s, Nt, N2t, X := galaxy.Cluster(ctx, events, verbose)

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
