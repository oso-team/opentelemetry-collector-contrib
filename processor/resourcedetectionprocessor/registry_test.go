// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

func TestDetectorRegistryReturnsCopy(t *testing.T) {
	registry := detectorRegistry()
	if len(registry) == 0 {
		t.Skip("registry is empty for this build tag configuration")
	}

	var detectorType internal.DetectorType
	for key := range registry {
		detectorType = key
		break
	}

	delete(registry, detectorType)

	assert.Contains(t, globalDetectorRegistry, detectorType)
}

func TestSortedDetectorTypes(t *testing.T) {
	registry := map[internal.DetectorType]internal.DetectorFactory{
		"system": nil,
		"env":    nil,
		"ec2":    nil,
	}

	assert.Equal(t, []string{"ec2", "env", "system"}, sortedDetectorTypes(registry))
}
