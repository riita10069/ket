package kubectl

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func (k *Kubectl) ApplyAllManifest(ctx context.Context, manifests []string, wait bool) error {
	if len(manifests) == 0 {
		return nil
	}

	var eg errgroup.Group
	for _, manifestPath := range manifests {
		manifestPath := manifestPath
		eg.Go(func() error {
			if wait {
				err := k.ApplyFileAndWait(ctx, manifestPath)
				if err != nil {
					return fmt.Errorf("failed to apply fixture and wait ready %s: %w", manifestPath, err)
				}
			} else {
				err := k.ApplyFile(ctx, manifestPath)
				if err != nil {
					return fmt.Errorf("failed to apply fixture %s: %w", manifestPath, err)
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (k *Kubectl) DeleteAllManifest(ctx context.Context, manifests []string, wait bool) error {
	if len(manifests) == 0 {
		return nil
	}

	var eg errgroup.Group
	for _, manifest := range manifests {
		manifest := manifest
		eg.Go(func() error {
			if wait {
				err := k.DeleteFileAndWait(ctx, manifest)
				if err != nil {
					return fmt.Errorf("failed to delete fixture because of applying %s: %w", manifest, err)
				}
			} else {
				err := k.DeleteFile(ctx, manifest)
				if err != nil {
					return fmt.Errorf("failed to delete fixture because of applying %s: %w", manifest, err)
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
