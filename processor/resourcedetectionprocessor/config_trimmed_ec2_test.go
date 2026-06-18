// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors && enable_ec2_detector

package resourcedetectionprocessor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ec2"
)

func TestTrimmedBuildLoadsEnabledDetectorConfig(t *testing.T) {
	cfg := NewFactory().CreateDefaultConfig()

	conf := confmap.NewFromStringMap(map[string]any{
		"detectors": []any{"ec2"},
		"ec2": map[string]any{
			"tags":         []any{"^tag1$", "^tag2$"},
			"max_attempts": 3,
			"max_backoff":  "20s",
		},
	})
	require.NoError(t, conf.Unmarshal(cfg))

	ec2Config, ok := cfg.(*Config).DetectorConfig.GetConfigFromType(ec2.TypeStr).(ec2.Config)
	require.True(t, ok)
	assert.Equal(t, []string{"^tag1$", "^tag2$"}, ec2Config.Tags)
	assert.Equal(t, 3, ec2Config.MaxAttempts)
	assert.Equal(t, 20*time.Second, ec2Config.MaxBackoff)
}
