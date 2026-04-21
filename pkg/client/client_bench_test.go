package client

import (
	"context"
	"testing"

	"digital.vasic.gandalfsolutions/pkg/types"
)

func BenchmarkSearchSolutions(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()
	ctx := context.Background()
	opts := types.SearchOptions{Query: "level", Limit: 50}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := c.SearchSolutions(ctx, opts); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetLevel(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer c.Close()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		if _, err := c.GetLevel(ctx, (i%8)+1); err != nil {
			b.Fatal(err)
		}
	}
}
