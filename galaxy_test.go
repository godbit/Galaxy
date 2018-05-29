package galaxy_test

import (
	"testing"

	galaxy "github.com/godbit/Galaxy"
)

func BenchmarkCluster(b *testing.B) {
	events, err := galaxy.ParseFile("testdata/data/ten.json")
	if err != nil {
		b.Errorf("%+v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		galaxy.Cluster(events)
	}
}
