package decorator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initDecoratorTest(_ context.Context) Decorator {
	return Decorator{}
}

func TestDecorator(t *testing.T) {
	ctx := context.Background()
	d := initDecoratorTest(ctx)

	_ = d.Init(ctx)

	defer func(d *Decorator, ctx context.Context) {
		_ = d.Deinit(ctx)
	}(&d, ctx)

	_, err := d.Call(ctx, "arg")
	assert.Equal(t, nil, err)
}
