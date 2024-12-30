package wechat

import (
	"context"
	"fmt"
)

const (
	name        = "wechat"
	description = "wechat tools"
)

type WeChat struct{}

func (w WeChat) Init(_ context.Context) error {
	return nil
}

func (w WeChat) Deinit(_ context.Context) error {
	return nil
}

func (w WeChat) Name(_ context.Context) string {
	return name
}

func (w WeChat) Description(_ context.Context) string {
	return description
}

func (w WeChat) Call(_ context.Context, args ...interface{}) (string, error) {
	return fmt.Sprintf("%v\n", args), nil
}
