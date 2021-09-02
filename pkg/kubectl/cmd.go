package kubectl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/riita10069/ket/pkg/util/slice"
	"k8s.io/apimachinery/pkg/types"
)

func (k *Kubectl) UseContext(ctx context.Context, clusterName string) error {
	args := []string{
		"config",
		"use-context",
		"kind-" + clusterName,
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to execute kubectl config use-context kind-%s: %w", clusterName, err)
	}

	return nil
}

func (k *Kubectl) ApplyKustomize(ctx context.Context, kustomizePath string) error {
	if kustomizePath == "" {
		return nil
	}
	args := []string{
		"apply",
		"-k",
		kustomizePath,
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to execute kubectl apply -k %s: %w", kustomizePath, err)
	}

	return nil
}

func (k *Kubectl) DeleteKustomize(ctx context.Context, kustomizePath string) error {
	args := []string{
		"delete",
		"-k",
		kustomizePath,
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to execute kubectl delete -k %s: %w", kustomizePath, err)
	}

	return nil
}

func (k *Kubectl) ApplyFile(ctx context.Context, filePath string) error {
	args := []string{
		"apply",
		"-f",
		filePath,
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to execute kubectl apply -f %s: %w", filePath, err)
	}

	return nil
}

func (k *Kubectl) DeleteFile(ctx context.Context, filePath string) error {
	args := []string{
		"delete",
		"-f",
		filePath,
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to execute kubectl delete -f %s: %w", filePath, err)
	}

	return nil
}

func (k *Kubectl) WaitFileForReady(ctx context.Context, filePath string) error {
	args := []string{
		"wait",
		"--filename",
		filePath,
		"--for",
		"condition=Ready",
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed or timeout to wait for applying %s: %w", filePath, err)
	}
	return nil
}

func (k *Kubectl) ApplyFileAndWait(ctx context.Context, filepath string) error {
	err := k.ApplyFile(ctx, filepath)
	if err != nil {
		return err
	}

	err = k.WaitFileForReady(ctx, filepath)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kubectl) DeleteFileAndWait(ctx context.Context, filePath string) error {
	args := []string{
		"delete",
		"-f",
		filePath,
		"--wait=true",
	}

	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to execute kubectl delete -f %s --wait=true: %w", filePath, err)
	}

	return nil
}

func (k *Kubectl) GetNamespacesList(ctx context.Context) ([]string, error) {
	kubectlArgs := []string{
		"get",
		"namespace",
		`-o=jsonpath='{.items[*].metadata.name}'`,
	}

	stdout, _, err := k.Capture(ctx, kubectlArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to capture kubectl get namespace: %w", err)
	}
	list := formatOutput(stdout)
	if list == "" {
		return nil, nil
	}
	return strings.Split(list, " "), nil
}

func (k *Kubectl) GetResourceNameList(ctx context.Context, namespace, resource string) ([]string, error) {
	args := []string{
		"get",
		resource,
		"-n",
		namespace,
		`-o=jsonpath='{.items[*].metadata.name}'`,
	}
	stdout, _, err := k.Capture(ctx, args)
	if err != nil {
		return nil, fmt.Errorf("failed to capture kubectl get %s -n %s: %w", resource, namespace, err)
	}
	list := formatOutput(stdout)
	if list == "" {
		return nil, nil
	}
	return strings.Split(list, " "), nil
}

func (k *Kubectl) GetResourceStatusList(ctx context.Context, namespacedName types.NamespacedName, resource string) (bool, error) {
	args := []string{
		"get",
		resource,
		namespacedName.Name,
		"-n",
		namespacedName.Namespace,
		`-o=jsonpath='-o=jsonpath='{.status.conditions[*].type}'`,
	}
	stdout, _, err := k.Capture(ctx, args)
	if err != nil {
		return false, fmt.Errorf("failed to get resorce status %s -n %s: %w", resource, namespacedName.Name, err)
	}
	list := formatOutput(stdout)
	if list == "" {
		return false, nil
	}
	statusList := strings.Split(list, " ")
	if slice.Contains([]string{"po", "pod", "pods"}, resource) {
		if !slice.Contains(statusList, "Ready") {
			return false, nil
		}
	} else if slice.Contains([]string{"deploy", "deployment", "deployments"}, resource) {
		if !slice.Contains(statusList, "Available") {
			return false, nil
		}
	}

	return true, nil
}

func (k *Kubectl) DeleteResource(ctx context.Context, name, namespace, resource string) error {
	args := []string{
		"delete",
		resource,
		name,
		"--namespace",
		namespace,
	}
	err := k.Execute(ctx, args)
	if err != nil {
		return fmt.Errorf("failed to capture kubectl delete %s -n %s: %w", resource, namespace, err)
	}
	return nil
}

// WaitAResource waits until deploy is ready.
func (k *Kubectl) WaitAResource(ctx context.Context, resource string, namespacedName types.NamespacedName) (ready bool, err error) {
	resource = strings.ToLower(resource)
	started := time.Now()
	checkInterval := 1 * time.Second
	timeout := 5 * time.Minute

	for {
		if time.Since(started) > timeout {
			return false, fmt.Errorf("waiting %s but it's time out", namespacedName.Name)
		}
		// First, Check if the resource exists.
		resourceNameList, err := k.GetResourceNameList(ctx, namespacedName.Namespace, resource)
		if err != nil {
			return false, fmt.Errorf("failed to get %s/%s in %s: %w", resource, namespacedName.Name, namespacedName.Namespace, err)
		}
		if slice.Contains(resourceNameList, namespacedName.Name) {
			// Second, Check whether the STATUS of the resource is READY or not.
			// However, check only for Pods, ReplicaSets, and Deployments.
			if slice.Contains([]string{"po", "pod", "pods", "deploy", "deployment", "deployments"}, resource) {
				ready, err := k.GetResourceStatusList(
					ctx,
					types.NamespacedName{
						Namespace: namespacedName.Namespace,
						Name:      namespacedName.Name,
					}, resource)
				if ready {
					return true, nil
				}
				if err != nil {
					return false, fmt.Errorf("failed to get current status: %w", err)
				}
			}
		}
		time.Sleep(checkInterval)
	}
}

func formatOutput(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
		if s[0] == '\'' && s[len(s)-1] == '\'' {
			return s[1 : len(s)-1]
		}
		if s[0] == '`' && s[len(s)-1] == '`' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
