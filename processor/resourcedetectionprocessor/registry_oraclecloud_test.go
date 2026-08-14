// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors && enable_resourcedetection_oraclecloud_detector

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

func TestOracleCloudDetectorRegistry(t *testing.T) {
	for detectorType, factory := range globalDetectorRegistry {
		require.NotNil(t, factory, "factory should not be nil for %q", detectorType)
	}

	assert.Contains(t, globalDetectorRegistry, internal.DetectorType("oraclecloud"))
}
