// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package containerapps // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/containerapps"

import containerappsconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/containerapps/config"

type Config = containerappsconfig.Config

func CreateDefaultConfig() Config {
	return containerappsconfig.CreateDefaultConfig()
}
