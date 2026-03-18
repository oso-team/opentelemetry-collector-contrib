// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oraclecloudauthextension

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "api_key_valid_default",
			cfg: Config{
				AuthType: AuthTypeAPIKey,
			},
		},
		{
			name: "api_key_profile_without_file",
			cfg: Config{
				AuthType: AuthTypeAPIKey,
				APIKey: APIKeyConfig{
					Profile: "DEFAULT",
				},
			},
			wantErr: true,
		},
		{
			name: "instance_principal_valid",
			cfg: Config{
				AuthType: AuthTypeInstancePrincipal,
			},
		},
		{
			name: "resource_principal_valid",
			cfg: Config{
				AuthType: AuthTypeResourcePrincipal,
			},
		},
		{
			name: "invalid_auth_type",
			cfg: Config{
				AuthType: "something_else",
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

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	assert.Equal(t, AuthTypeAPIKey, cfg.AuthType)
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}
