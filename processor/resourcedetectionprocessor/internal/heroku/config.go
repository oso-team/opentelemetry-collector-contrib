// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package heroku // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/heroku"

import herokuconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/heroku/config"

type Config = herokuconfig.Config

func CreateDefaultConfig() Config {
	return herokuconfig.CreateDefaultConfig()
}
