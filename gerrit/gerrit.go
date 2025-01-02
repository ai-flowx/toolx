package gerrit

import (
	"context"
	"regexp"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/pkg/errors"
)

const (
	name        = "gerrit"
	description = "gerrit tools"

	remote = "https://android.googlesource.com/"
	branch = "main"

	fromLen = 3
	fileLen = 4
)

type Gerrit struct {
	Project string
	Branch  string
	Patch   Patch
}

type Patch struct {
	Hash    string
	Author  string
	Email   string
	Date    string
	Subject string
	File    []File
	Diff    map[string]string
}

type File struct {
	Name      string
	Insertion int
	Deletion  int
}

func (g *Gerrit) Init(_ context.Context) error {
	return nil
}

func (g *Gerrit) Deinit(_ context.Context) error {
	return nil
}

func (g *Gerrit) Name(_ context.Context) string {
	return name
}

func (g *Gerrit) Description(_ context.Context) string {
	return description
}

func (g *Gerrit) Call(ctx context.Context, _ func(context.Context, interface{}) (interface{}, error), args ...interface{}) (string, error) {
	if len(args) == 0 || args[0].(string) == "" {
		return "", errors.New("invalid arguments\n")
	}

	if err := g.load(ctx, args[0].(string)); err != nil {
		return "", err
	}

	if err := g.clone(ctx); err != nil {
		return "", err
	}

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.clean(ctx)
	}(g, ctx)

	if err := g.config(ctx); err != nil {
		return "", err
	}

	if err := g.commit(ctx); err != nil {
		return "", err
	}

	url, err := g.push(ctx)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (g *Gerrit) load(ctx context.Context, content string) error {
	changes, summary, err := gitdiff.Parse(strings.NewReader(content))
	if err != nil {
		return errors.Wrap(err, "failed to parse diff\n")
	}

	if err := g.parseSummary(ctx, summary); err != nil {
		return errors.Wrap(err, "failed to parse summary\n")
	}

	if err := g.parseChange(ctx, changes); err != nil {
		return errors.Wrap(err, "failed to parse change\n")
	}

	return nil
}

func (g *Gerrit) clone(_ context.Context) error {
	return nil
}

func (g *Gerrit) config(_ context.Context) error {
	return nil
}

func (g *Gerrit) commit(_ context.Context) error {
	return nil
}

func (g *Gerrit) push(_ context.Context) (string, error) {
	return "", nil
}

func (g *Gerrit) clean(_ context.Context) error {
	return nil
}

func (g *Gerrit) parseSummary(_ context.Context, content string) error {
	lines := strings.Split(content, "\n")

	if len(lines) == 0 {
		return errors.New("invalid content\n")
	}

	if parts := strings.Fields(lines[0]); len(parts) > 0 {
		g.Patch.Hash = parts[0]
	}

	authorRegexp := regexp.MustCompile(`From: ([^<]+)<([^>]+)>`)
	for _, line := range lines {
		if strings.HasPrefix(line, "From:") {
			matches := authorRegexp.FindStringSubmatch(line)
			if len(matches) == fromLen {
				g.Patch.Author = strings.TrimSpace(matches[1])
				g.Patch.Email = strings.TrimSpace(matches[2])
			}
		}
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "Date:") {
			g.Patch.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		}
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "Subject:") {
			g.Patch.Subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		}
	}

	fileRegexp := regexp.MustCompile(`(\S+)\s+\|\s+(\d+)\s+([+-]+)`)
	for _, line := range lines {
		matches := fileRegexp.FindStringSubmatch(line)
		if len(matches) == fileLen {
			file := File{
				Name:      matches[1],
				Insertion: strings.Count(matches[3], "+"),
				Deletion:  strings.Count(matches[3], "-"),
			}
			g.Patch.File = append(g.Patch.File, file)
		}
	}

	return nil
}

func (g *Gerrit) parseChange(_ context.Context, content []*gitdiff.File) error {
	return nil
}
