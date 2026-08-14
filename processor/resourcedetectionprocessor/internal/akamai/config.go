// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package akamai // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/akamai"

import akamaiconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/akamai/config"

type Config = akamaiconfig.Config

func CreateDefaultConfig() Config {
	return akamaiconfig.CreateDefaultConfig()
}
