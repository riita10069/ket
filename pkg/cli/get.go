package cli

import (
	"context"
	"fmt"
	"os"
)

func Get(ctx context.Context, cli CLI) error {
	_, err := os.Stat(cli.Path())
	if err == nil {
		return nil
	}
	if err := get(ctx, cli); err != nil {
		return fmt.Errorf("failed to download and install %s: %w", cli.Name(), err)
	}
	return nil
}
