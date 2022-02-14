package kubectl

import (
	"context"
	"fmt"
	"github.com/riita10069/ket/pkg/cli"
	"os"
	"path/filepath"
	"runtime"
)

type Kubectl struct {
	name           string
	version        string
	binDir         string
	url            string
	kubeConfigPath string
}

func NewKubectl(version, binDir, kubeConfigFilePath string) *Kubectl {
	return &Kubectl{
		version:        version,
		name:           "kubectl",
		binDir:         binDir,
		url:            fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/v%s/bin/%s/%s/kubectl", version, runtime.GOOS, runtime.GOARCH),
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

func (k *Kubectl) URL() string {
	return k.url
}

func (k *Kubectl) Envs() []string {
	return []string{
		"KUBECONFIG=" + k.kubeConfigPath,
	}
}

// Execute If OutPut is necessary, use Capture. Execute uses os.Stderr.
func (k *Kubectl) Execute(ctx context.Context, args []string) error {
	return cli.Run(ctx, k, args, os.Stdout, os.Stderr)
}

// Capture execute command with returning outs as string.
func (k *Kubectl) Capture(ctx context.Context, args []string) (stdout string, stderr string, err error) {
	return cli.Capture(ctx, k, args)
}
