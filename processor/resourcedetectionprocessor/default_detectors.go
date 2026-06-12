// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !remove_all_resourcedetection_detectors

package resourcedetectionprocessor

import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/env"

func defaultDetectors() []string {
	return []string{env.TypeStr}
}
