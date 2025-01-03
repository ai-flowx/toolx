package gerrit

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

const (
	name        = "gerrit"
	description = "gerrit tools"
	count       = 3

	remoteUrl    = "https://android.googlesource.com"
	remoteBranch = "main"

	hashCol = 2
	fromCol = 3
	fileCol = 4

	userName  = "name"
	userEmail = "name@example.com"

	reviewUrl = `https://[-a-zA-Z0-9]+\.googlesource\.com/c/[^/]+/[^/]+/[^/]+/\+/\d+`
)

type Gerrit struct {
	Path    string
	Project string
	Branch  string
	Patch   Patch

	repo *git.Repository
}

type Patch struct {
	Hash    string
	Author  string
	Email   string
	Date    string
	Subject string
	File    []File
	Diff    map[string]*gitdiff.File
}

type File struct {
	Name      string
	Insertion int
	Deletion  int
}

func (g *Gerrit) Init(_ context.Context) error {
	g.Path = ""
	g.Project = ""
	g.Branch = remoteBranch
	g.Patch = Patch{
		File: []File{},
		Diff: map[string]*gitdiff.File{},
	}

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
	if len(args) == 0 {
		return "", errors.New("invalid arguments\n")
	}

	_ = g.Init(ctx)

	_patch := args[0].(string)

	if len(args) == count-1 {
		g.Project = args[1].(string)
	} else if len(args) >= count {
		g.Project = args[1].(string)
		g.Branch = args[2].(string)
	}

	if err := g.load(ctx, _patch); err != nil {
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

	if err := g.apply(ctx); err != nil {
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

func (g *Gerrit) clone(ctx context.Context) error {
	var err error

	if g.Project == "" || g.Branch == "" {
		return errors.New("invalid arguments\n")
	}

	g.Path = path.Join(os.TempDir(), g.Project)
	if _, err = os.Stat(g.Path); err == nil {
		return errors.Wrap(err, "path already exists\n")
	}

	if g.repo, err = git.PlainCloneContext(ctx, g.Path, false, &git.CloneOptions{
		URL:             fmt.Sprintf("%s/%s", remoteUrl, g.Project),
		ReferenceName:   plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", g.Branch)),
		SingleBranch:    true,
		Depth:           1,
		InsecureSkipTLS: true,
	}); err != nil {
		return errors.Wrap(err, "failed to clone repo\n")
	}

	return nil
}

func (g *Gerrit) config(_ context.Context) error {
	cfg, err := g.repo.Config()
	if err != nil {
		return errors.Wrap(err, "failed to get config\n")
	}

	cfg.User.Name = userName
	cfg.User.Email = userEmail

	if err := g.repo.SetConfig(cfg); err != nil {
		return errors.Wrap(err, "failed to set config\n")
	}

	return nil
}

func (g *Gerrit) apply(_ context.Context) error {
	for _, item := range g.Patch.File {
		p := path.Join(g.Path, item.Name)
		f, err := os.Open(p)
		if err != nil {
			return errors.Wrap(err, "failed to open file\n")
		}
		var buf bytes.Buffer
		if err := gitdiff.Apply(&buf, f, g.Patch.Diff[item.Name]); err != nil {
			_ = f.Close()
			return errors.Wrap(err, "failed to apply patch\n")
		}
		s, err := os.Stat(p)
		if err != nil {
			_ = f.Close()
			return errors.Wrap(err, "failed to stat file\n")
		}
		if err := os.WriteFile(p, buf.Bytes(), s.Mode().Perm()); err != nil {
			_ = f.Close()
			return errors.Wrap(err, "failed to write file\n")
		}
		_ = f.Close()
	}

	return nil
}

func (g *Gerrit) commit(_ context.Context) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree\n")
	}

	for _, item := range g.Patch.File {
		if _, err := wt.Add(item.Name); err != nil {
			return errors.Wrap(err, "failed to add file\n")
		}
	}

	_, err = wt.Commit(g.Patch.Subject, &git.CommitOptions{
		Author: &object.Signature{
			Name:  g.Patch.Author,
			Email: g.Patch.Email,
		},
	})

	if err != nil {
		return errors.Wrap(err, "failed to commit change\n")
	}

	return nil
}

func (g *Gerrit) push(ctx context.Context) (string, error) {
	var out io.Writer

	defer func(g *Gerrit, ctx context.Context) {
		_ = g.reset(ctx)
	}(g, ctx)

	if err := g.repo.Push(&git.PushOptions{
		Progress:        out,
		InsecureSkipTLS: true,
	}); err != nil {
		return "", errors.Wrap(err, "failed to push change\n")
	}

	url, err := g.parsePush(ctx, fmt.Sprint(out))
	if err != nil {
		return "", errors.Wrap(err, "failed to parse push\n")
	}

	return url, nil
}

func (g *Gerrit) reset(_ context.Context) error {
	wt, err := g.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "failed to get worktree\n")
	}

	ref, err := g.repo.Head()
	if err != nil {
		return errors.Wrap(err, "failed to get head\n")
	}

	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return errors.Wrap(err, "failed to get commit\n")
	}

	parentCommit, err := commit.Parent(0)
	if err != nil {
		return errors.Wrap(err, "failed to get parent\n")
	}

	if err = wt.Reset(&git.ResetOptions{
		Commit: parentCommit.Hash,
		Mode:   git.HardReset,
	}); err != nil {
		return errors.Wrap(err, "failed to reset commit\n")
	}

	return nil
}

func (g *Gerrit) clean(_ context.Context) error {
	_ = os.RemoveAll(g.Path)
	_ = os.Remove(g.Path)

	return nil
}

func (g *Gerrit) parseSummary(_ context.Context, content string) error {
	lines := strings.Split(content, "\n")

	if len(lines) == 0 {
		return errors.New("invalid content\n")
	}

	if parts := strings.Fields(lines[0]); len(parts) >= hashCol {
		g.Patch.Hash = parts[1]
	}

	authorRegexp := regexp.MustCompile(`From: ([^<]+)<([^>]+)>`)
	for _, line := range lines {
		if strings.HasPrefix(line, "From:") {
			matches := authorRegexp.FindStringSubmatch(line)
			if len(matches) == fromCol {
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
		if len(matches) == fileCol {
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
	if len(content) == 0 || len(g.Patch.File) != len(content) {
		return errors.New("invalid content\n")
	}

	for index, item := range g.Patch.File {
		g.Patch.Diff[item.Name] = content[index]
	}

	return nil
}

func (g *Gerrit) parsePush(_ context.Context, content string) (string, error) {
	var match string

	pattern := regexp.MustCompile(reviewUrl)

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if match = pattern.FindString(line); match != "" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return match, nil
}
