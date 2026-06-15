// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors && enable_oraclecloud_detector

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

func TestOracleCloudDetectorRegistry(t *testing.T) {
	expectedTypes := []internal.DetectorType{
		"oraclecloud",
	}

	registry := detectorRegistry()
	actualTypes := make([]internal.DetectorType, 0, len(registry))
	for detectorType, factory := range registry {
		require.NotNil(t, factory, "factory should not be nil for %q", detectorType)
		actualTypes = append(actualTypes, detectorType)
	}

	assert.ElementsMatch(t, expectedTypes, actualTypes)
}

func TestOracleCloudDetectorConfigRegistry(t *testing.T) {
	config := detectorCreateDefaultConfig()

	assert.NotNil(t, config.GetConfigFromType("oraclecloud"))
	assert.Nil(t, config.GetConfigFromType("env"))
}
