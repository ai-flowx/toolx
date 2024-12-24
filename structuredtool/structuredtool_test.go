package structuredtool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initStructuredToolTest(_ context.Context) StructuredTool {
	return StructuredTool{}
}

func TestStructuredTool(t *testing.T) {
	ctx := context.Background()
	d := initStructuredToolTest(ctx)

	_ = d.Init(ctx)

	defer func(d *StructuredTool, ctx context.Context) {
		_ = d.Deinit(ctx)
	}(&d, ctx)

	_, err := d.Call(ctx, "arg")
	assert.Equal(t, nil, err)
}
