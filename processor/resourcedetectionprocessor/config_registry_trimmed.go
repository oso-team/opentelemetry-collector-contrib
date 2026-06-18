// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors

package resourcedetectionprocessor

import (
	"go.opentelemetry.io/collector/confmap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

type detectorConfigRegistration struct {
	configKey string
	create    func() internal.DetectorConfig
	unmarshal func(*confmap.Conf) (internal.DetectorConfig, error)
}

var detectorConfigRegistrations = map[internal.DetectorType]detectorConfigRegistration{}

func registerDetectorConfig[T any](detectorType internal.DetectorType, createDefaultConfig func() T) {
	registerDetectorConfigWithKey(detectorType, string(detectorType), createDefaultConfig)
}

func registerDetectorConfigWithKey[T any](detectorType internal.DetectorType, configKey string, createDefaultConfig func() T) {
	detectorConfigRegistrations[detectorType] = detectorConfigRegistration{
		configKey: configKey,
		create: func() internal.DetectorConfig {
			return createDefaultConfig()
		},
		unmarshal: func(conf *confmap.Conf) (internal.DetectorConfig, error) {
			config := createDefaultConfig()
			if conf.ToStringMap() == nil {
				return config, nil
			}
			if err := conf.Unmarshal(&config); err != nil {
				return nil, err
			}
			return config, nil
		},
	}
}

func detectorCreateDefaultConfig() DetectorConfig {
	config := DetectorConfig{configs: map[internal.DetectorType]internal.DetectorConfig{}}
	for detectorType, registration := range detectorConfigRegistrations {
		config.configs[detectorType] = registration.create()
	}
	return config
}

func (d *DetectorConfig) GetConfigFromType(detectorType internal.DetectorType) internal.DetectorConfig {
	if d == nil {
		return nil
	}
	return d.configs[detectorType]
}
