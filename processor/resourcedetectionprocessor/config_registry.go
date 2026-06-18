// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor

import (
	"reflect"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

type detectorConfigGetter func(*DetectorConfig) any
type detectorConfigDefault func(*DetectorConfig)

var (
	detectorConfigGetters  = map[internal.DetectorType]detectorConfigGetter{}
	detectorConfigKeys     = map[string]struct{}{}
	detectorConfigDefaults []detectorConfigDefault
)

func registerDetectorConfig[T any](detectorType internal.DetectorType, createDefaultConfig func() T) {
	registerDetectorConfigWithKey(detectorType, string(detectorType), createDefaultConfig)
}

func registerDetectorConfigWithKey[T any](detectorType internal.DetectorType, configKey string, createDefaultConfig func() T) {
	fieldIndex := detectorConfigFieldIndex(configKey)
	fieldType := reflect.TypeOf(DetectorConfig{}).Field(fieldIndex).Type
	defaultType := reflect.TypeFor[T]()
	if !defaultType.AssignableTo(fieldType) {
		panic("resource detector config " + configKey + " has default config type " + defaultType.String() + ", want " + fieldType.String())
	}

	detectorConfigKeys[configKey] = struct{}{}
	detectorConfigGetters[detectorType] = func(config *DetectorConfig) any {
		return reflect.ValueOf(config).Elem().Field(fieldIndex).Interface()
	}
	detectorConfigDefaults = append(detectorConfigDefaults, func(config *DetectorConfig) {
		reflect.ValueOf(config).Elem().Field(fieldIndex).Set(reflect.ValueOf(createDefaultConfig()))
	})
}

func detectorConfigFieldIndex(configKey string) int {
	configType := reflect.TypeOf(DetectorConfig{})
	for i := 0; i < configType.NumField(); i++ {
		tag := configType.Field(i).Tag.Get("mapstructure")
		if strings.Split(tag, ",")[0] == configKey {
			return i
		}
	}
	panic("resource detector config key " + configKey + " is not defined on DetectorConfig")
}

func detectorCreateDefaultConfig() DetectorConfig {
	var config DetectorConfig
	for _, defaultConfig := range detectorConfigDefaults {
		defaultConfig(&config)
	}
	return config
}

func (d *DetectorConfig) GetConfigFromType(detectorType internal.DetectorType) internal.DetectorConfig {
	getter, ok := detectorConfigGetters[detectorType]
	if !ok {
		return nil
	}
	return internal.DetectorConfig(getter(d))
}
