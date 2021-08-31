package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func get(ctx context.Context, cli CLI) error {
	// ensure install dir
	if err := os.MkdirAll(cli.Dir(), 0o755); err != nil {
		return fmt.Errorf("can't create %s for %s dir: %w", cli.Dir(), cli.Name(), err)
	}

	// ensure download dir
	tmpDir, err := ioutil.TempDir(cli.Dir(), "cache")
	if err != nil {
		return fmt.Errorf("can't create cache directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	err = download(ctx, cli, tmpDir)
	if err != nil {
		return fmt.Errorf("can't download %s: %w", cli.Name(), err)
	}

	err = install(ctx, cli, tmpDir)
	if err != nil {
		return fmt.Errorf("can't install %s: %w", cli.Name(), err)
	}

	return nil
}

func download(ctx context.Context, cli CLI, tmpDir string) error {
	downloadPath := filepath.Join(tmpDir, cli.Name()+"-"+cli.Version())

	out, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("can't create download path: %w", err)
	}
	defer out.Close()

	resp, err := cli.Download(ctx)
	if err != nil {
		return fmt.Errorf("can't download %s: %w", cli.Name(), err)
	}
	defer resp.Close()

	_, err = io.Copy(out, resp)
	if err != nil {
		return fmt.Errorf("can't write the downloaded file: %w", err)
	}
	return nil
}

func install(ctx context.Context, cli CLI, tmpDir string) error {
	downloadPath := filepath.Join(tmpDir, cli.Name()+"-"+cli.Version())

	downloadedFile, err := os.Open(downloadPath)
	if err != nil {
		return fmt.Errorf("can't open downloadedFile file: %s: %w", downloadPath, err)
	}

	err = cli.Install(ctx, downloadedFile)
	if err != nil {
		return fmt.Errorf("can't extract executable binary for %s: %w", cli.Name(), err)
	}

	return nil
}
