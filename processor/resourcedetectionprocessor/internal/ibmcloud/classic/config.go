// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package classic // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/ibmcloud/classic"

import classicconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/ibmcloud/classic/config"

type Config = classicconfig.Config

func CreateDefaultConfig() Config {
	return classicconfig.CreateDefaultConfig()
}
