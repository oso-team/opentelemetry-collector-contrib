// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor

import (
	"context"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/processor"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

type metadataLifecycleDetector struct{}

func (*metadataLifecycleDetector) Detect(context.Context) (pcommon.Resource, string, error) {
	return pcommon.NewResource(), "", nil
}

func init() {
	registerDetector("metadata_test", func(processor.Settings, internal.DetectorConfig, bool) (internal.Detector, error) {
		return &metadataLifecycleDetector{}, nil
	})
}
