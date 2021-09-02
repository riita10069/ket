package skaffold

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
}

func NewSkaffold(version, binDir, kubeConfigPath string) *Skaffold {
	return &Skaffold{
		version:        version,
		name:           "skaffold",
		binDir:         binDir,
		kubeConfigPath: kubeConfigPath,
	}
}

func (s *Skaffold) Version() string {
	return s.version
}

func (s *Skaffold) Name() string {
	return s.name
}

func (s *Skaffold) Path() string {
	return filepath.Join(s.binDir, s.name)
}

func (s *Skaffold) Dir() string {
	return s.binDir
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

func (s *Skaffold) DownloadURL() string {
	return fmt.Sprintf("https://storage.googleapis.com/skaffold/releases/v%s/skaffold-%s-%s", s.Version(), runtime.GOOS, runtime.GOARCH)
}

func (s *Skaffold) Download(ctx context.Context) (io.ReadCloser, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, s.DownloadURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize request: %w", err)
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("skaffold url responses error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("skaffold url responses bad status: %s", resp.Status)
	}
	return resp.Body, nil
}

func (s *Skaffold) Install(ctx context.Context, file io.Reader) error {
	bin, err := os.Create(s.Path())
	if err != nil {
		return fmt.Errorf("can't create skaffold binary path: %w", err)
	}
	defer bin.Close()

	_, err = io.Copy(bin, file)
	if err != nil {
		return fmt.Errorf("failed to copy skaffold binary to %s: %w", s.Path(), err)
	}

	err = os.Chmod(s.Path(), 0o755)
	if err != nil {
		return fmt.Errorf("failed to chmod when skaffold binary path: %w", err)
	}
	return nil
}

// Execute If OutPut is necessary, use Capture. Execute uses os.Stderr.
func (s *Skaffold) Execute(ctx context.Context, args []string) error {
	return cli.Run(ctx, s, args, os.Stdout, os.Stderr)
}

// Capture execute command with returning outs as string.
func (s *Skaffold) Capture(ctx context.Context, args []string) (stdout string, stderr string, err error) {
	return cli.Capture(ctx, s, args)
}
