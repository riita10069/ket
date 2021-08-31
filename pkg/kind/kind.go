package kind

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

type Kind struct {
	version           string
	kubernetesVersion string
	name              string
	binDir            string
	kubeConfigPath    string
}

func NewKind(kindVersion, kubernetesVersion, binDir, kubeConfigPath string) *Kind {
	return &Kind{
		version:           kindVersion,
		name:              "kind",
		binDir:            binDir,
		kubeConfigPath:    kubeConfigPath,
		kubernetesVersion: kubernetesVersion,
	}
}

func (k *Kind) Version() string {
	return k.version
}

func (k *Kind) Name() string {
	return k.name
}

func (k *Kind) Path() string {
	return filepath.Join(k.binDir, k.name)
}

func (k *Kind) Dir() string {
	return k.binDir
}

func (k *Kind) Envs() []string {
	return []string{}
}

func (k *Kind) DownloadURL() string {
	return fmt.Sprintf("https://github.com/kubernetes-sigs/kind/releases/download/v%s/kind-%s-%s", k.Version(), runtime.GOOS, runtime.GOARCH)
}

func (k *Kind) Download(ctx context.Context) (io.ReadCloser, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, k.DownloadURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize request: %w", err)
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kind url responses error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kind url responses bad status: %s", resp.Status)
	}
	return resp.Body, nil
}

func (k *Kind) Install(ctx context.Context, file io.Reader) error {
	bin, err := os.Create(k.Path())
	if err != nil {
		return fmt.Errorf("can't create kind binary path: %w", err)
	}
	defer bin.Close()

	_, err = io.Copy(bin, file)
	if err != nil {
		return fmt.Errorf("failed to copy kind binary to %s: %w", k.Path(), err)
	}

	err = os.Chmod(k.Path(), 0o755)
	if err != nil {
		return fmt.Errorf("failed to chmod when kind binary path: %w", err)
	}
	return nil
}

// Execute If OutPut is necessary, use Capture. Execute uses os.Stderr.
func (k *Kind) Execute(ctx context.Context, args []string) error {
	return cli.Run(ctx, k, args, os.Stdout, os.Stderr)
}

// Capture execute command with returning outs as string
func (k *Kind) Capture(ctx context.Context, args []string) (stdout string, stderr string, err error) {
	return cli.Capture(ctx, k, args)
}
