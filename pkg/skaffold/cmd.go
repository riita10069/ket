package skaffold

import (
	"context"
	"time"
)

// Run exec skaffold run -f {filename}
// if logs is true, it outputs kubectl logs to stdout.
func (s *Skaffold) Run(ctx context.Context, filename string, logs bool) error {
	args := []string{
		"dev",
		"-f",
		filename,
		"--port-forward",
	}

	if logs {
		args = append(args, "--tail")
	}

	go func(ctx context.Context) {
		if err := s.Execute(ctx, args); err != nil {
			// fmt.Printf("failed to exec %v\n build or deploy resource of %s: %v", args, filename, err)
			// FIXME: We should not allow this goroutine to cause a selfish panic.
			time.Sleep(10 * time.Millisecond)
			panic(err)
		}
	}(ctx)
	return nil
}
