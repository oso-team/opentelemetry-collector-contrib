// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package docker // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/docker"

import dockerconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/docker/config"

type Config = dockerconfig.Config

func CreateDefaultConfig() Config {
	return dockerconfig.CreateDefaultConfig()
}
