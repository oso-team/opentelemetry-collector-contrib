// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oraclecloud // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/oraclecloud"

import oraclecloudconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/oraclecloud/config"

type Config = oraclecloudconfig.Config

func CreateDefaultConfig() Config {
	return oraclecloudconfig.CreateDefaultConfig()
}
