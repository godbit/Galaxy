//+build cgo

package galaxy

import (
	"context"
	"fmt"
	"log"

	"github.com/karlek/progress/barcli"
)

func inner(ctx context.Context, imin, imax int, events []Event, bar *barcli.Bar, verbose bool, c chan Result) {
	var result Result

	fmt.Println("inner loop in C")

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
