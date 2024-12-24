package runnables

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initRunnablesTest(_ context.Context) Runnables {
	return Runnables{}
}

func TestRunnables(t *testing.T) {
	ctx := context.Background()
	d := initRunnablesTest(ctx)

	_ = d.Init(ctx)

	defer func(d *Runnables, ctx context.Context) {
		_ = d.Deinit(ctx)
	}(&d, ctx)

	_, err := d.Call(ctx, "arg")
	assert.Equal(t, nil, err)
}
