// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package lambda // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/lambda"

import lambdaconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/lambda/config"

type Config = lambdaconfig.Config

func CreateDefaultConfig() Config {
	return lambdaconfig.CreateDefaultConfig()
}
