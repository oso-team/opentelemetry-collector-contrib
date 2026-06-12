// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveAllDetectorsCreateDefaultConfigHasNoDetectors(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	assert.Empty(t, cfg.Detectors)
}
