package kubectl

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

type Kubectl struct {
	version        string
	name           string
	binDir         string
	kubeConfigPath string
}

func NewKubectl(version, binDir, kubeConfigFilePath string) *Kubectl {
	return &Kubectl{
		version:        version,
		name:           "kubectl",
		binDir:         binDir,
		kubeConfigPath: kubeConfigFilePath,
	}
}

func (k *Kubectl) Version() string {
	return k.version
}

func (k *Kubectl) Name() string {
	return k.name
}

func (k *Kubectl) Path() string {
	return filepath.Join(k.binDir, k.name)
}

func (k *Kubectl) Dir() string {
	return k.binDir
}

func (k *Kubectl) Envs() []string {
	return []string{
		"KUBECONFIG=" + k.kubeConfigPath,
	}
}

func (k *Kubectl) DownloadURL() string {
	return fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/v%s/bin/%s/%s/kubectl", k.Version(), runtime.GOOS, runtime.GOARCH)
}

func (k *Kubectl) Download(ctx context.Context) (io.ReadCloser, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, k.DownloadURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize request: %w", err)
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("kubectl url responses error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kubectl url responses bad status: %s", resp.Status)
	}
	return resp.Body, nil
}

func (k *Kubectl) Install(ctx context.Context, file io.Reader) error {
	bin, err := os.Create(k.Path())
	if err != nil {
		return fmt.Errorf("can't create Kubectl binary path: %w", err)
	}
	defer bin.Close()

	_, err = io.Copy(bin, file)
	if err != nil {
		return fmt.Errorf("failed to copy Kubectl binary to %s: %w", k.Path(), err)
	}

	err = os.Chmod(k.Path(), 0o755)
	if err != nil {
		return fmt.Errorf("failed to chmod when Kubectl binary path: %w", err)
	}
	return nil
}

// Execute If OutPut is necessary, use Capture. Execute uses os.Stderr.
func (k *Kubectl) Execute(ctx context.Context, args []string) error {
	return cli.Run(ctx, k, args, os.Stdout, os.Stderr)
}

// Capture execute command with returning outs as string
func (k *Kubectl) Capture(ctx context.Context, args []string) (stdout string, stderr string, err error) {
	return cli.Capture(ctx, k, args)
}
