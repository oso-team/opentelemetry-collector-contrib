package oraclecloudmonitoringexporter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configoptional"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid_with_auth_extension",
			cfg: Config{
				Region: "us-phoenix-1",
				Auth: configoptional.Some(configauth.Config{
					AuthenticatorID: component.MustNewID("oraclecloudauth"),
				}),
			},
		},
		{
			name:    "missing_region",
			cfg:     Config{},
			wantErr: true,
		},
		{
			name: "missing_auth",
			cfg: Config{
				Region: "us-phoenix-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
