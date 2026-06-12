// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !remove_all_resourcedetection_detectors

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBuildCreateDefaultConfigUsesEnvDetector(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	assert.Equal(t, []string{"env"}, cfg.Detectors)
}
