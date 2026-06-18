// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

import (
	"time"

	"go.opentelemetry.io/collector/config/confighttp"
)

// Config defines configuration for Resource processor.
type Config struct {
	// Detectors is an ordered list of named detectors that should be
	// run to attempt to detect resource information.
	Detectors []string `mapstructure:"detectors"`
	// Override indicates whether any existing resource attributes
	// should be overridden or preserved. Defaults to true.
	Override bool `mapstructure:"override"`
	// DetectorConfig is a list of settings specific to all detectors
	DetectorConfig DetectorConfig `mapstructure:",squash"`
	// HTTP client settings for the detector
	// Timeout default is 5s
	confighttp.ClientConfig `mapstructure:",squash"`
	// If > 0, periodically re-run detection for all configured detectors.
	// When 0 (default), no periodic refresh occurs.
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
}
