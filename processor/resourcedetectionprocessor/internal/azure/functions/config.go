// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package functions // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/functions"

import functionsconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/functions/config"

type Config = functionsconfig.Config

func CreateDefaultConfig() Config {
	return functionsconfig.CreateDefaultConfig()
}
