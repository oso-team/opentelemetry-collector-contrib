// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors && enable_oraclecloud_detector

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
)

func TestOracleCloudOnlyBuildRejectsDisabledDetectorConfigBlock(t *testing.T) {
	cfg := NewFactory().CreateDefaultConfig()

	conf := confmap.NewFromStringMap(map[string]any{
		"detectors": []any{"oraclecloud"},
		"ec2": map[string]any{
			"tags": []any{"^tag1$"},
		},
	})

	require.ErrorContains(t, conf.Unmarshal(cfg), `detector config "ec2" is not registered in this build`)
}

func TestOracleCloudOnlyBuildAllowsEnabledDetectorConfigBlock(t *testing.T) {
	cfg := NewFactory().CreateDefaultConfig()

	conf := confmap.NewFromStringMap(map[string]any{
		"detectors":   []any{"oraclecloud"},
		"oraclecloud": map[string]any{},
	})

	require.NoError(t, conf.Unmarshal(cfg))
}
