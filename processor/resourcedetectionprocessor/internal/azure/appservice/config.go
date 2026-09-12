// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package appservice // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/appservice"

import appserviceconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/appservice/config"

type Config = appserviceconfig.Config

func CreateDefaultConfig() Config {
	return appserviceconfig.CreateDefaultConfig()
}
