// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oraclecloudauthextension // import "github.com/open-telemetry/opentelemetry-collector-contrib/extension/oraclecloudauthextension"

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
)

const (
	AuthTypeAPIKey            = "api_key"
	AuthTypeInstancePrincipal = "instance_principal"
	AuthTypeResourcePrincipal = "resource_principal"
)

var _ component.Config = (*Config)(nil)

type Config struct {
	AuthType string `mapstructure:"auth_type"`

	APIKey            APIKeyConfig            `mapstructure:"api_key"`
	InstancePrincipal InstancePrincipalConfig `mapstructure:"instance_principal"`
	ResourcePrincipal ResourcePrincipalConfig `mapstructure:"resource_principal"`
}

type APIKeyConfig struct {
	ConfigFile           string `mapstructure:"config_file"`
	Profile              string `mapstructure:"profile"`
	PrivateKeyPassphrase string `mapstructure:"private_key_passphrase"`
}

type InstancePrincipalConfig struct {
	Region string `mapstructure:"region"`
}

type ResourcePrincipalConfig struct {
	Region string `mapstructure:"region"`
}

func (cfg *Config) Validate() error {
	switch cfg.AuthType {
	case AuthTypeAPIKey:
		if cfg.APIKey.Profile != "" && cfg.APIKey.ConfigFile == "" {
			return errors.New(`"api_key.profile" requires "api_key.config_file"`)
		}
		return nil
	case AuthTypeInstancePrincipal:
		return nil
	case AuthTypeResourcePrincipal:
		return nil
	default:
		return fmt.Errorf(`"auth_type" must be one of [%s, %s, %s]`, AuthTypeAPIKey, AuthTypeInstancePrincipal, AuthTypeResourcePrincipal)
	}
}
