package cli

import (
	"context"
	"io"
)

type CLI interface {
	Version() string
	Name() string
	Path() string
	Dir() string
	Download(context.Context) (io.ReadCloser, error)
	Install(context.Context, io.Reader) error
	Envs() []string
}
