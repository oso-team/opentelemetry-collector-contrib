// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

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
