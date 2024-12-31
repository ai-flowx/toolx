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

	c := func(context.Context, interface{}) (interface{}, error) {
		return nil, nil
	}

	_, err := d.Call(ctx, c, "arg")
	assert.Equal(t, nil, err)
}
