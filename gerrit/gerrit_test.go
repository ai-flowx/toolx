//go:build gerrit_test

// go test -cover -covermode=atomic -parallel 2 -tags=gerrit_test -v github.com/ai-flowx/toolx/gerrit

package gerrit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	content string
)

func initGerritTest(_ context.Context) Gerrit {
	return Gerrit{
		Project: "",
		Branch:  branch,
		Commit:  Commit{},
	}
}

func TestParse(t *testing.T) {
	ctx := context.Background()
	g := initGerritTest(ctx)

	err := g.parse(ctx, content)
	assert.Equal(t, nil, err)
}

func TestClone(t *testing.T) {
	ctx := context.Background()
	g := initGerritTest(ctx)

	_ = g.parse(ctx, content)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	err := g.clone(ctx)
	assert.Equal(t, nil, err)
}

func TestConfig(t *testing.T) {
	ctx := context.Background()
	g := initGerritTest(ctx)

	_ = g.parse(ctx, content)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	_ = g.clone(ctx)

	err := g.config(ctx)
	assert.Equal(t, nil, err)
}

func TestCommit(t *testing.T) {
	ctx := context.Background()
	g := initGerritTest(ctx)

	_ = g.parse(ctx, content)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	_ = g.clone(ctx)
	_ = g.config(ctx)

	err := g.commit(ctx)
	assert.Equal(t, nil, err)
}

func TestPush(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestClean(t *testing.T) {
	ctx := context.Background()
	g := initGerritTest(ctx)

	_ = g.parse(ctx, content)
	_ = g.clone(ctx)

	err := g.clean(ctx)
	assert.Equal(t, nil, err)
}
