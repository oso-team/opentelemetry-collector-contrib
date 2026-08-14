// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sapi // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/k8sapi"

import k8sapiconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/k8sapi/config"

type Config = k8sapiconfig.Config

func CreateDefaultConfig() Config {
	return k8sapiconfig.CreateDefaultConfig()
}
