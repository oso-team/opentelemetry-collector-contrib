// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

func TestSortedDetectorTypes(t *testing.T) {
	registry := map[internal.DetectorType]internal.DetectorFactory{
		"system": nil,
		"env":    nil,
		"ec2":    nil,
	}

	assert.Equal(t, []string{"ec2", "env", "system"}, sortedDetectorTypes(registry))
}
