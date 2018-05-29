// Package knox implements the Knox test for space-time interaction.
//
// References:
//
// Kulldorff, Martin, and Ulf Hjalmars. "The Knox method and other tests  for
// space‐time interaction." Biometrics 55.2 (1999): 544-552.
//
// https://pdfs.semanticscholar.org/91fb/c7143da6efcdca6964d603f454c69c3911c3.pdf
package knox

import "math"

// Test performs the Knox test and returns the number of pairs, expected value,
// and variance.
//
// Let n be the total number of cases, so that N = n(n-1)/2 distinct pair of
// cases.
//
// Let Nt be the number of case pairs that are closer to each other in time as
// compared to some specified temporal distance.
//
// Likewise, let Ns be the number of pairs close in space as defined by some
// geographic distance.
//
// Finally, let X be the number of case pairs that are close both in time and
// space.
//
// N2s is the number of pairs of `case pairs close in space` that have one case
// in common and where N2t is defined equivalently for time.
func Test(Ns, N2s, Nt, N2t, X, n float64) (N, E, V float64) {
	// Number of pairs
	N = n * (n - 1) / 2

	// Expected value
	E = Nt * Ns / N

	// Variance
	V = Ns*Nt/N + 4*N2s*N2t/(n*(n-1)*(n-2)) +
		4*(Ns*(Ns-1)-N2s)*(Nt*(Nt-1)-N2t)/(n*(n-1)*(n-2)*(n-3)) -
		math.Pow(Ns*Nt/N, 2)

	return N, E, V
}
