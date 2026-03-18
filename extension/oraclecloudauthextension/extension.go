// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oraclecloudauthextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/oraclecloudauthextension"

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

var (
	newAPIKeyProvider = func(cfg APIKeyConfig) (common.ConfigurationProvider, error) {
		switch {
		case cfg.ConfigFile == "":
			return common.DefaultConfigProvider(), nil
		case cfg.Profile == "":
			return common.ConfigurationProviderFromFile(cfg.ConfigFile, cfg.PrivateKeyPassphrase)
		default:
			return common.ConfigurationProviderFromFileWithProfile(cfg.ConfigFile, cfg.Profile, cfg.PrivateKeyPassphrase)
		}
	}
	newInstancePrincipalProvider = func(cfg InstancePrincipalConfig) (common.ConfigurationProvider, error) {
		if cfg.Region == "" {
			return auth.InstancePrincipalConfigurationProvider()
		}
		return auth.InstancePrincipalConfigurationProviderForRegion(common.StringToRegion(cfg.Region))
	}
	newResourcePrincipalProvider = func(cfg ResourcePrincipalConfig) (common.ConfigurationProvider, error) {
		if cfg.Region == "" {
			return auth.ResourcePrincipalConfigurationProvider()
		}
		return auth.ResourcePrincipalConfigurationProviderForRegion(common.StringToRegion(cfg.Region))
	}
)

// ConfigurationProvider returns OCI SDK configuration provider for downstream components.
type ConfigurationProvider interface {
	ConfigurationProvider() common.ConfigurationProvider
}

type oracleCloudAuth struct {
	component.StartFunc
	component.ShutdownFunc
	provider common.ConfigurationProvider
}

var (
	_ extension.Extension   = (*oracleCloudAuth)(nil)
	_ ConfigurationProvider = (*oracleCloudAuth)(nil)
)

func newOracleCloudAuthExtension(cfg *Config) (*oracleCloudAuth, error) {
	provider, err := configProviderFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &oracleCloudAuth{
		provider: provider,
	}, nil
}

func (oa *oracleCloudAuth) ConfigurationProvider() common.ConfigurationProvider {
	return oa.provider
}

func configProviderFromConfig(cfg *Config) (common.ConfigurationProvider, error) {
	var (
		provider common.ConfigurationProvider
		err      error
	)

	switch cfg.AuthType {
	case AuthTypeAPIKey:
		provider, err = newAPIKeyProvider(cfg.APIKey)
	case AuthTypeInstancePrincipal:
		provider, err = newInstancePrincipalProvider(cfg.InstancePrincipal)
	case AuthTypeResourcePrincipal:
		provider, err = newResourcePrincipalProvider(cfg.ResourcePrincipal)
	default:
		return nil, fmt.Errorf("unsupported auth_type %q", cfg.AuthType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OCI configuration provider for auth_type %q: %w", cfg.AuthType, err)
	}

	return provider, nil
}

func (*oracleCloudAuth) Start(context.Context, component.Host) error {
	return nil
}

func (*oracleCloudAuth) Shutdown(context.Context) error {
	return nil
}
