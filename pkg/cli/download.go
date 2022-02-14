package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func get(ctx context.Context, cli CLI) error {
	if err := os.MkdirAll(cli.Dir(), 0o755); err != nil {
		return fmt.Errorf("can't create %s for %s dir: %w", cli.Dir(), cli.Name(), err)
	}

	out, err := os.Create(cli.Path())
	if err != nil {
		return fmt.Errorf("can't create download path: %w", err)
	}
	defer out.Close()

	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, cli.URL(), nil)
	if err != nil {
		return fmt.Errorf("failed to initialize request: %w", err)
	}
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("url responses error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("url responses bad status: %s", resp.Status)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("can't write the downloaded file: %w", err)
	}

	err = os.Chmod(cli.Path(), 0o755)
	if err != nil {
		return fmt.Errorf("failed to chmod when kind binary path: %w", err)
	}

	return nil
}
