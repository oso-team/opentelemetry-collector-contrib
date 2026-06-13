// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors && enable_oraclecloud_detector

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/metadata"
)

// This file tests the behavior of a "trimmed build" in which go "build tags" are used to exclude/include
// specific resource detection processor modules from the build.  As such, this test is intended to be run
// differently than the standard unit tests.
//
// go test -tags='remove_all_resourcedetection_detectors,enable_oraclecloud_detector' . -run 'TestOracleCloudOnlyBuildCreatesConfiguredProcessor|TestOracleCloudOnlyBuildRejectsDisabledDetector' -count=1
//
// Note also the "go:build" directive at the top of the file, which appears as a comment but is actually a compile-time directive
// which means that this test should only run if we've used build tags to exclude all resource detection processors,
// and selectively enable the oracle cloud processor
func TestOracleCloudOnlyBuildCreatesConfiguredProcessor(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)
	oCfg.Detectors = []string{"oraclecloud"}

	tp, err := factory.CreateTraces(t.Context(), processortest.NewNopSettings(metadata.Type), cfg, consumertest.NewNop())
	require.NoError(t, err)
	assert.NotNil(t, tp)
}

// A detector that was excluded by build tags must fail at processor creation
// instead of being silently accepted from config.
func TestOracleCloudOnlyBuildRejectsDisabledDetector(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)
	oCfg.Detectors = []string{"env"}

	tp, err := factory.CreateTraces(t.Context(), processortest.NewNopSettings(metadata.Type), cfg, consumertest.NewNop())
	require.ErrorContains(t, err, "invalid detector key: env")
	assert.Nil(t, tp)
}
