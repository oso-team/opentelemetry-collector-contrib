// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ec2 // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ec2"

import ec2config "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ec2/config"

type Config = ec2config.Config

func CreateDefaultConfig() Config {
	return ec2config.CreateDefaultConfig()
}
