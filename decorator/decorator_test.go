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

	c := func(context.Context, interface{}) (interface{}, error) {
		return nil, nil
	}

	_, err := d.Call(ctx, c, "arg")
	assert.Equal(t, nil, err)
}
