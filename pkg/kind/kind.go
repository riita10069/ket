package kind

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/riita10069/ket/pkg/cli"
)

type Kind struct {
	name              string
	version           string
	kubernetesVersion string
	binDir            string
	url               string
	kubeConfigPath    string
}

func NewKind(kindVersion, kubernetesVersion, binDir, kubeConfigPath string) *Kind {
	return &Kind{
		version:           kindVersion,
		name:              "kind",
		binDir:            binDir,
		url:               fmt.Sprintf("https://github.com/kubernetes-sigs/kind/releases/download/v%s/kind-%s-%s", kindVersion, runtime.GOOS, runtime.GOARCH),
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

func (k *Kind) URL() string {
	return k.url
}

func (k *Kind) Envs() []string {
	return []string{}
}

// Execute If OutPut is necessary, use Capture. Execute uses os.Stderr.
func (k *Kind) Execute(ctx context.Context, args []string) error {
	return cli.Run(ctx, k, args, os.Stdout, os.Stderr)
}

// Capture execute command with returning outs as string.
func (k *Kind) Capture(ctx context.Context, args []string) (stdout string, stderr string, err error) {
	return cli.Capture(ctx, k, args)
}
