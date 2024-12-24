package basetool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initBaseToolTest(_ context.Context) BaseTool {
	return BaseTool{}
}

func TestBaseTool(t *testing.T) {
	ctx := context.Background()
	d := initBaseToolTest(ctx)

	_ = d.Init(ctx)

	defer func(d *BaseTool, ctx context.Context) {
		_ = d.Deinit(ctx)
	}(&d, ctx)

	_, err := d.Call(ctx, "arg")
	assert.Equal(t, nil, err)
}
