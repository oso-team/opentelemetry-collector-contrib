// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package openshift // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openshift"

import openshiftconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openshift/config"

type Config = openshiftconfig.Config

func CreateDefaultConfig() Config {
	return openshiftconfig.CreateDefaultConfig()
}
