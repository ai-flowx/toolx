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

const (
	pushGerritTest = `
Enumerating objects: 5, done.
Counting objects: 100% (5/5), done.
Delta compression using up to 2 threads
Compressing objects: 100% (3/3), done.
Writing objects: 100% (3/3), 353 bytes | 353.00 KiB/s, done.
Total 3 (delta 2), reused 0 (delta 0), pack-reused 0
remote: Resolving deltas: 100% (2/2)
remote: Waiting for private key checker: 1/1 objects left
remote: Update redirected to refs/for/refs/heads/main.
remote: Processing changes: refs: 1, new: 1, done
remote:
remote: SUCCESS
remote:
remote:   https://android-review.googlesource.com/c/platform/build/soong/+/3435262 Test only [NEW]
remote:
To https://android.googlesource.com/platform/build/soong
 * [new reference]       HEAD -> refs/for/master`
)

func initGerritTest(_ context.Context) (Gerrit, string) {
	g := Gerrit{
		Project: "platform/build/soong",
		Branch:  remoteBranch,
		Patch: Patch{
			File: []File{},
			Diff: map[string]*gitdiff.File{},
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

func TestHook(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	_ = g.clone(ctx)
	_ = g.config(ctx)

	err := g.hook(ctx)
	assert.Equal(t, nil, err)
}

func TestApply(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	_ = g.clone(ctx)
	_ = g.config(ctx)
	_ = g.hook(ctx)

	err := g.apply(ctx)
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
	_ = g.hook(ctx)
	_ = g.apply(ctx)

	err := g.commit(ctx)
	assert.Equal(t, nil, err)
}

func TestPush(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	_ = g.clone(ctx)
	_ = g.config(ctx)
	_ = g.hook(ctx)
	_ = g.apply(ctx)
	_ = g.commit(ctx)

	err := g.push(ctx)
	assert.Equal(t, nil, err)
}

func TestReset(t *testing.T) {
	ctx := context.Background()
	g, d := initGerritTest(ctx)

	_ = g.load(ctx, d)

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(&g, ctx)

	_ = g.clone(ctx)
	_ = g.config(ctx)
	_ = g.hook(ctx)
	_ = g.apply(ctx)
	_ = g.commit(ctx)

	err := g.reset(ctx)
	assert.Equal(t, nil, err)
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

func TestParsePush(t *testing.T) {
	ctx := context.Background()
	g, _ := initGerritTest(ctx)

	_, err := g.parsePush(ctx, pushGerritTest)
	assert.Equal(t, nil, err)
}
