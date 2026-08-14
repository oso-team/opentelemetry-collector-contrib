// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aks // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/aks"

import aksconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/aks/config"

type Config = aksconfig.Config

func CreateDefaultConfig() Config {
	return aksconfig.CreateDefaultConfig()
}
