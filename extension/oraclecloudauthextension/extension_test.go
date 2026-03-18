// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oraclecloudauthextension

import (
	"errors"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigProviderFromConfig(t *testing.T) {
	origAPI := newAPIKeyProvider
	origInstance := newInstancePrincipalProvider
	origResource := newResourcePrincipalProvider
	t.Cleanup(func() {
		newAPIKeyProvider = origAPI
		newInstancePrincipalProvider = origInstance
		newResourcePrincipalProvider = origResource
	})

	tests := []struct {
		name      string
		cfg       Config
		setup     func()
		wantErr   bool
		assertion func(t *testing.T, provider common.ConfigurationProvider)
	}{
		{
			name: "api_key",
			cfg:  Config{AuthType: AuthTypeAPIKey},
			setup: func() {
				newAPIKeyProvider = func(cfg APIKeyConfig) (common.ConfigurationProvider, error) {
					assert.Equal(t, APIKeyConfig{}, cfg)
					return common.DefaultConfigProvider(), nil
				}
			},
		},
		{
			name: "instance_principal",
			cfg: Config{
				AuthType: AuthTypeInstancePrincipal,
				InstancePrincipal: InstancePrincipalConfig{
					Region: "us-ashburn-1",
				},
			},
			setup: func() {
				newInstancePrincipalProvider = func(cfg InstancePrincipalConfig) (common.ConfigurationProvider, error) {
					assert.Equal(t, "us-ashburn-1", cfg.Region)
					return common.DefaultConfigProvider(), nil
				}
			},
		},
		{
			name: "resource_principal",
			cfg: Config{
				AuthType: AuthTypeResourcePrincipal,
				ResourcePrincipal: ResourcePrincipalConfig{
					Region: "us-phoenix-1",
				},
			},
			setup: func() {
				newResourcePrincipalProvider = func(cfg ResourcePrincipalConfig) (common.ConfigurationProvider, error) {
					assert.Equal(t, "us-phoenix-1", cfg.Region)
					return common.DefaultConfigProvider(), nil
				}
			},
		},
		{
			name:    "invalid_auth_type",
			cfg:     Config{AuthType: "invalid"},
			wantErr: true,
		},
		{
			name: "provider_error",
			cfg:  Config{AuthType: AuthTypeAPIKey},
			setup: func() {
				newAPIKeyProvider = func(APIKeyConfig) (common.ConfigurationProvider, error) {
					return nil, errors.New("boom")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			provider, err := configProviderFromConfig(&tt.cfg)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, provider)
			if tt.assertion != nil {
				tt.assertion(t, provider)
			}
		})
	}
}

func TestNewOracleCloudAuthExtension(t *testing.T) {
	origAPI := newAPIKeyProvider
	t.Cleanup(func() {
		newAPIKeyProvider = origAPI
	})
	newAPIKeyProvider = func(APIKeyConfig) (common.ConfigurationProvider, error) {
		return common.DefaultConfigProvider(), nil
	}

	ext, err := newOracleCloudAuthExtension(&Config{AuthType: AuthTypeAPIKey})
	require.NoError(t, err)
	require.NotNil(t, ext)
	require.NotNil(t, ext.ConfigurationProvider())
}
