package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Run(ctx context.Context, cli CLI, args []string, stdout, stderr io.Writer) error {
	err := Get(ctx, cli)
	if err != nil {
		return fmt.Errorf("failed to ensure %s: %w", cli.Name(), err)
	}
	cmd := exec.CommandContext(ctx, cli.Path(), args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, cli.Envs()...)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to exec %v %s %v: %w", cmd.Env, cmd.Path, cmd.Args, err)
	}
	return nil
}

func Capture(ctx context.Context, cli CLI, args []string) (string, string, error) {
	outb := new(bytes.Buffer)
	errb := new(bytes.Buffer)

	err := Run(ctx, cli, args, outb, errb)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute command %s %v: %w", cli.Name(), args, err)
	}

	return outb.String(), errb.String(), err
}
