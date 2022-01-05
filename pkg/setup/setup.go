package setup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/riita10069/ket/pkg/k8s"
	"github.com/riita10069/ket/pkg/kind"
	"github.com/riita10069/ket/pkg/kubectl"
	"github.com/riita10069/ket/pkg/skaffold"
)

type Option func(*KET) error

func WithBinaryDirectory(binDir string) Option {
	return func(k *KET) error {
		k.binDir = binDir
		return nil
	}
}

func WithKindVersion(kindVersion string) Option {
	return func(k *KET) error {
		k.kindVersion = kindVersion
		return nil
	}
}

func WithKindClusterName(kindClusterName string) Option {
	return func(k *KET) error {
		k.kindClusterName = kindClusterName
		return nil
	}
}

func WithKubernetesVersion(kubernetesVersion string) Option {
	return func(k *KET) error {
		k.kubernetesVersion = kubernetesVersion
		return nil
	}
}

func WithKubeconfigPath(kubeconfigPath string) Option {
	return func(k *KET) error {
		k.kubeconfigPath = kubeconfigPath
		return nil
	}
}

func WithNotCRD() Option {
	return func(k *KET) error {
		k.isThereCRD = false
		return nil
	}
}

func WithCRDKustomizePath(crdKustomizePath string) Option {
	return func(k *KET) error {
		k.crdKustomizePath = crdKustomizePath
		return nil
	}
}

func WithUseSkaffold() Option {
	return func(k *KET) error {
		k.useSkaffold = true
		return nil
	}
}

func WithSkaffoldVersion(skaffoldVersion string) Option {
	return func(k *KET) error {
		k.skaffoldVersion = skaffoldVersion
		return nil
	}
}

func WithSkaffoldYaml(skaffoldYaml string) Option {
	return func(k *KET) error {
		k.skaffoldYaml = skaffoldYaml
		return nil
	}
}

type KET struct {
	binDir            string
	kindVersion       string
	kindClusterName   string
	kubernetesVersion string
	kubeconfigPath    string
	isThereCRD        bool
	crdKustomizePath  string
	useSkaffold       bool
	skaffoldVersion   string
	skaffoldYaml      string
}

func NewKET() *KET {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	return &KET{
		binDir:            "./bin",
		kindVersion:       "0.11.0",
		kindClusterName:   "ket",
		kubernetesVersion: "1.20.2",
		kubeconfigPath:    filepath.Join(homeDir, ".kube", "config"),
		isThereCRD:        true,
		crdKustomizePath:  "",
		useSkaffold:       false,
		skaffoldVersion:   "1.26.1",
		skaffoldYaml:      "./skaffold/skaffold.yaml",
	}
}

type ClientSet struct {
	ClientGo *k8s.ClientGo
	Kubectl  *kubectl.Kubectl
	Kind     *kind.Kind
	Skaffold *skaffold.Skaffold
}

func Start(ctx context.Context, options ...Option) (*ClientSet, error) {
	ket := NewKET()
	for _, option := range options {
		err := option(ket)
		if err != nil {
			return nil, fmt.Errorf("failed to run options: %w", err)
		}
	}

	cliSet := &ClientSet{}
	kind := kind.NewKind(ket.kindVersion, ket.kubernetesVersion, ket.binDir, ket.kubeconfigPath)
	cliSet.Kind = kind

	err := kind.DeleteCluster(ctx, ket.kindClusterName)
	if err != nil {
		return nil, fmt.Errorf("failed to delete kind cluster %s: %w", ket.kindClusterName, err)
	}

	err = kind.CreateCluster(ctx, ket.kindClusterName)
	if err != nil {
		return nil, fmt.Errorf("failed to create kind cluster %s: %w", ket.kindClusterName, err)
	}

	clientGo, err := k8s.NewClientGo(ket.kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create client-go: %w", err)
	}
	cliSet.ClientGo = clientGo

	kubectl := kubectl.NewKubectl(ket.kubernetesVersion, ket.binDir, ket.kubeconfigPath)
	cliSet.Kubectl = kubectl

	err = kubectl.UseContext(ctx, ket.kindClusterName)
	if err != nil {
		return nil, fmt.Errorf("failed to use context: %w", err)
	}

	if ket.isThereCRD {
		err = kubectl.ApplyKustomize(ctx, ket.crdKustomizePath)
		if err != nil {
			return nil, fmt.Errorf("failed to apply crd yaml: %w", err)
		}
	}

	// TODO
	// Waiting for resources to be applied by kustomize Just before.
	// It should be guaranteed that the resource is created.
	time.Sleep(3 * time.Second)

	if ket.useSkaffold {
		skaffold := skaffold.NewSkaffold(ket.skaffoldVersion, ket.binDir, ket.kubeconfigPath)
		cliSet.Skaffold = skaffold
		err = skaffold.Run(ctx, ket.skaffoldYaml, false)
		if err != nil {
			return nil, fmt.Errorf("failed to skaffold run: %w", err)
		}
	}

	return cliSet, nil
}
