// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor

import (
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

var globalDetectorRegistry = map[internal.DetectorType]internal.DetectorFactory{}

func registerDetector(detectorType internal.DetectorType, factory internal.DetectorFactory) {
	if _, ok := globalDetectorRegistry[detectorType]; ok {
		panic(fmt.Sprintf("duplicate detector registration for %q", detectorType))
	}

	globalDetectorRegistry[detectorType] = factory
}

func detectorRegistry() map[internal.DetectorType]internal.DetectorFactory {
	registry := make(map[internal.DetectorType]internal.DetectorFactory, len(globalDetectorRegistry))
	for detectorType, factory := range globalDetectorRegistry {
		registry[detectorType] = factory
	}
	return registry
}