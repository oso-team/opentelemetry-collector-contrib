// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package upcloud // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/upcloud"

import upcloudconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/upcloud/config"

type Config = upcloudconfig.Config

func CreateDefaultConfig() Config {
	return upcloudconfig.CreateDefaultConfig()
}
