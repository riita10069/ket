package kind

import (
	"context"
	"fmt"
)

func (k *Kind) CreateCluster(ctx context.Context, clusterName string) error {
	args := []string{
		"create",
		"cluster",
		"--name",
		clusterName,
		"--image",
		"kindest/node:v" + k.kubernetesVersion,
		"--kubeconfig",
		k.kubeConfigPath,
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to create kind cluster: %w", err)
	}
	return nil
}

func (k *Kind) DeleteCluster(ctx context.Context, clusterName string) error {
	args := []string{
		"delete",
		"cluster",
		"--name",
		clusterName,
		"--kubeconfig",
		k.kubeConfigPath,
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to delete kind cluster: %w", err)
	}
	return nil
}
