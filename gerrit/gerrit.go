package gerrit

import (
	"context"

	"github.com/pkg/errors"
)

const (
	name        = "gerrit"
	description = "gerrit tools"
)

type Gerrit struct{}

func (g Gerrit) Init(_ context.Context) error {
	return nil
}

func (g Gerrit) Deinit(_ context.Context) error {
	return nil
}

func (g Gerrit) Name(_ context.Context) string {
	return name
}

func (g Gerrit) Description(_ context.Context) string {
	return description
}

func (g Gerrit) Call(ctx context.Context, _ func(context.Context, interface{}) (interface{}, error), args ...interface{}) (string, error) {
	if len(args) == 0 {
		return "", errors.New("invalid arguments\n")
	}

	if err := g.clone(ctx); err != nil {
		return "", err
	}

	if err := g.patch(ctx); err != nil {
		return "", err
	}

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

func (g Gerrit) clone(_ context.Context) error {
	return nil
}

func (g Gerrit) patch(_ context.Context) error {
	return nil
}

func (g Gerrit) config(_ context.Context) error {
	return nil
}

func (g Gerrit) commit(_ context.Context) error {
	return nil
}

func (g Gerrit) push(_ context.Context) (string, error) {
	return "", nil
}
