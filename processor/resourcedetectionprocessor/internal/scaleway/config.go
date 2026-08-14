// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package scaleway // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/scaleway"

import scalewayconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/scaleway/config"

type Config = scalewayconfig.Config

func CreateDefaultConfig() Config {
	return scalewayconfig.CreateDefaultConfig()
}
