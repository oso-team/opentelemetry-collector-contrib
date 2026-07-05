// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package config // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm/config"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm/internal/metadata"
)

// APIConfig mirrors k8sconfig.APIConfig without importing client-go.
type APIConfig struct {
	AuthType     string  `mapstructure:"auth_type"`
	Context      string  `mapstructure:"context"`
	KubeAPIQPS   float32 `mapstructure:"kube_api_qps"`
	KubeAPIBurst int     `mapstructure:"kube_api_burst"`
}

type Config struct {
	APIConfig          `mapstructure:",squash"`
	ResourceAttributes metadata.ResourceAttributesConfig `mapstructure:"resource_attributes"`
}

func CreateDefaultConfig() Config {
	return Config{
		APIConfig:          APIConfig{AuthType: "serviceAccount"},
		ResourceAttributes: metadata.DefaultResourceAttributesConfig(),
	}
}
