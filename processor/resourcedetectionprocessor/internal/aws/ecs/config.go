// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ecs // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ecs"

import ecsconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ecs/config"

type Config = ecsconfig.Config

func CreateDefaultConfig() Config {
	return ecsconfig.CreateDefaultConfig()
}
