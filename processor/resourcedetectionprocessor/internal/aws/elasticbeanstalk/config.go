// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package elasticbeanstalk // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/elasticbeanstalk"

import elasticbeanstalkconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/elasticbeanstalk/config"

type Config = elasticbeanstalkconfig.Config

func CreateDefaultConfig() Config {
	return elasticbeanstalkconfig.CreateDefaultConfig()
}
