// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors

package resourcedetectionprocessor

import (
	"fmt"
	"reflect"
	"strings"

	"go.opentelemetry.io/collector/confmap"
)

func (c *Config) Unmarshal(conf *confmap.Conf) error {
	if err := rejectUnregisteredDetectorConfigBlocks(conf); err != nil {
		return err
	}

	type rawConfig Config
	return conf.Unmarshal((*rawConfig)(c))
}

func rejectUnregisteredDetectorConfigBlocks(conf *confmap.Conf) error {
	configType := reflect.TypeOf(DetectorConfig{})
	for i := 0; i < configType.NumField(); i++ {
		configKey := strings.Split(configType.Field(i).Tag.Get("mapstructure"), ",")[0]
		if configKey == "" || !conf.IsSet(configKey) {
			continue
		}
		if _, ok := detectorConfigKeys[configKey]; ok {
			continue
		}
		return fmt.Errorf("detector config %q is not registered in this build", configKey)
	}
	return nil
}
