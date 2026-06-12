// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build remove_all_resourcedetection_detectors

package resourcedetectionprocessor

func defaultDetectors() []string {
	return nil
}
