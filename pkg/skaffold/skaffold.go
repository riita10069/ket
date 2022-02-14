package skaffold

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/riita10069/ket/pkg/cli"
)

type Skaffold struct {
	version        string
	name           string
	binDir         string
	kubeConfigPath string
	url            string
}

func NewSkaffold(version, binDir, kubeConfigPath string) *Skaffold {
	return &Skaffold{
		version:        version,
		name:           "skaffold",
		binDir:         binDir,
		kubeConfigPath: kubeConfigPath,
		url:            fmt.Sprintf("https://storage.googleapis.com/skaffold/releases/v%s/skaffold-%s-%s", version, runtime.GOOS, runtime.GOARCH),
	}
}

func (s *Skaffold) Name() string {
	return s.name
}

func (s *Skaffold) Version() string {
	return s.version
}

func (s *Skaffold) Path() string {
	return filepath.Join(s.binDir, s.name)
}

func (s *Skaffold) Dir() string {
	return s.binDir
}

func (s *Skaffold) URL() string {
	return s.url
}

func (s *Skaffold) Envs() []string {
	pwd, err := os.Getwd()
	if err != nil {
		return []string{}
	}
	binDir := filepath.Join(pwd, s.binDir)

	return []string{
		"PATH=" + binDir + ":" + os.ExpandEnv("$PATH"),
		"KUBECONFIG=" + s.kubeConfigPath,
	}
}

// Execute If OutPut is necessary, use Capture. Execute uses os.Stderr.
func (s *Skaffold) Execute(ctx context.Context, args []string) error {
	return cli.Run(ctx, s, args, os.Stdout, os.Stderr)
}

// Capture execute command with returning outs as string.
func (s *Skaffold) Capture(ctx context.Context, args []string) (stdout string, stderr string, err error) {
	return cli.Capture(ctx, s, args)
}
