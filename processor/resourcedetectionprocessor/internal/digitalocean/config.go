// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package digitalocean // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/digitalocean"

import digitaloceanconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/digitalocean/config"

type Config = digitaloceanconfig.Config

func CreateDefaultConfig() Config {
	return digitaloceanconfig.CreateDefaultConfig()
}
