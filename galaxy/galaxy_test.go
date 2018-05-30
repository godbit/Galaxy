package galaxy_test

import (
	"context"
	"testing"

	"github.com/godbit/Galaxy/galaxy"
)

func BenchmarkCluster(b *testing.B) {
	events, err := galaxy.ParseFile("testdata/data/month.json")
	if err != nil {
		b.Errorf("%+v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		galaxy.Cluster(context.TODO(), events, false)
	}
}
