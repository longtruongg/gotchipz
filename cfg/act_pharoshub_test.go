package cfg_test

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/ethclient"
	"main/cfg"
	"testing"
)

func TestSignaturePharosHub(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		provider *ethclient.Client
		prik     *ecdsa.PrivateKey
		want     interface{}
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := cfg.SignaturePharosHub(context.Background(), tt.provider, tt.prik)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SignaturePharosHub() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SignaturePharosHub() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("SignaturePharosHub() = %v, want %v", got, tt.want)
			}
		})
	}
}
