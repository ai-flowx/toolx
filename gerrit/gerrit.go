package gerrit

import (
	"context"
	"fmt"
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

func (g Gerrit) Call(_ context.Context, args ...interface{}) (string, error) {
	return fmt.Sprintf("%v\n", args), nil
}
