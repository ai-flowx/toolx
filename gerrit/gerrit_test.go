//go:build gerrit_test

// go test -cover -covermode=atomic -parallel 2 -tags=gerrit_test -v github.com/ai-flowx/toolx/gerrit

package gerrit

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/stretchr/testify/assert"
)

func initGerritTest(_ context.Context) (Gerrit, string) {
	g := Gerrit{
		Project: "",
		Branch:  branch,
		Patch: Patch{
			File: []File{},
			Diff: map[string]string{},
		},
	}

	d, _ := os.ReadFile("../test/gerrit/test.patch")

	return g, string(d)
}

func TestLoad(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	err := g.load(ctx, d)
	assert.Equal(t, nil, err)
}

func TestClone(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	err := g.clone(ctx)
	assert.Equal(t, nil, err)
}

func TestConfig(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	_ = g.clone(ctx)

	err := g.config(ctx)
	assert.Equal(t, nil, err)
}

func TestCommit(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)

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
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)
	_ = g.clone(ctx)

	err := g.clean(ctx)
	assert.Equal(t, nil, err)
}

func TestParseSummary(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_, s, _ := gitdiff.Parse(strings.NewReader(d))

	err := g.parseSummary(ctx, s)
	assert.Equal(t, nil, err)
}

func TestParseChange(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	c, s, _ := gitdiff.Parse(strings.NewReader(d))
	_ = g.parseSummary(ctx, s)

	err := g.parseChange(ctx, c)
	assert.Equal(t, nil, err)
}
