package cli_test

import (
	"context"
	"testing"

	"github.com/riita10069/ket/pkg/cli"
	"github.com/riita10069/ket/pkg/kind"
	"github.com/riita10069/ket/pkg/kubectl"
	"github.com/riita10069/ket/pkg/skaffold"
)

const (
	binDir         = "./bin"
	kubeconfigPath = "./kubeconfig"
)

func Test_get(t *testing.T) {
	type args struct {
		cli cli.CLI
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get kubectl binary",
			args: args{
				kubectl.NewKubectl("1.21.2", binDir, kubeconfigPath),
			},
			wantErr: false,
		},
		{
			name: "get skaffold binary",
			args: args{
				skaffold.NewSkaffold("1.26.1", binDir, kubeconfigPath),
			},
			wantErr: false,
		},
		{
			name: "get kind binary",
			args: args{
				kind.NewKind("0.11.0", "1.20.2", binDir, kubeconfigPath),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if err := cli.Get(context.Background(), tt.args.cli); (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	//err := os.RemoveAll(binDir)
	//if err != nil {
	//	t.Errorf("failed to remove binary directory: %v", err)
	//	return
	//}
}
