// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package nova // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openstack/nova"

import novaconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openstack/nova/config"

type Config = novaconfig.Config

func CreateDefaultConfig() Config {
	return novaconfig.CreateDefaultConfig()
}
