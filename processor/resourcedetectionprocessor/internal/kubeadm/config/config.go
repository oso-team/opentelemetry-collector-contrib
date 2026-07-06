// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package config // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm/config"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig/k8sconfigtypes"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm/internal/metadata"
)

type Config struct {
	k8sconfigtypes.APIConfig `mapstructure:",squash"`
	ResourceAttributes       metadata.ResourceAttributesConfig `mapstructure:"resource_attributes"`
}

func CreateDefaultConfig() Config {
	return Config{
		APIConfig:          k8sconfigtypes.APIConfig{AuthType: k8sconfigtypes.AuthTypeServiceAccount},
		ResourceAttributes: metadata.DefaultResourceAttributesConfig(),
	}
}
