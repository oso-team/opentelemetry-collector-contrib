// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors

package resourcedetectionprocessor

import (
	"time"

	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/confmap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

// DetectorConfig contains user-specified configurations unique to enabled detectors.
type DetectorConfig struct {
	configs map[internal.DetectorType]internal.DetectorConfig
}

func (c *Config) Unmarshal(conf *confmap.Conf) error {
	detectorConfig := detectorCreateDefaultConfig()
	remaining := conf.ToStringMap()
	for detectorType, registration := range detectorConfigRegistrations {
		if !conf.IsSet(registration.configKey) {
			continue
		}
		sub, err := conf.Sub(registration.configKey)
		if err != nil {
			return err
		}
		config, err := registration.unmarshal(sub)
		if err != nil {
			return err
		}
		detectorConfig.configs[detectorType] = config
		delete(remaining, registration.configKey)
	}

	withoutDetectorConfig := configWithoutDetectorConfig{
		Detectors:       c.Detectors,
		Override:        c.Override,
		ClientConfig:    c.ClientConfig,
		RefreshInterval: c.RefreshInterval,
	}
	if err := confmap.NewFromStringMap(remaining).Unmarshal(&withoutDetectorConfig); err != nil {
		return err
	}

	c.Detectors = withoutDetectorConfig.Detectors
	c.Override = withoutDetectorConfig.Override
	c.ClientConfig = withoutDetectorConfig.ClientConfig
	c.RefreshInterval = withoutDetectorConfig.RefreshInterval
	c.DetectorConfig = detectorConfig
	return nil
}

type configWithoutDetectorConfig struct {
	// Detectors is an ordered list of named detectors that should be
	// run to attempt to detect resource information.
	Detectors []string `mapstructure:"detectors"`
	// Override indicates whether any existing resource attributes
	// should be overridden or preserved. Defaults to true.
	Override bool `mapstructure:"override"`
	// HTTP client settings for the detector
	// Timeout default is 5s
	confighttp.ClientConfig `mapstructure:",squash"`
	// If > 0, periodically re-run detection for all configured detectors.
	// When 0 (default), no periodic refresh occurs.
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
}
