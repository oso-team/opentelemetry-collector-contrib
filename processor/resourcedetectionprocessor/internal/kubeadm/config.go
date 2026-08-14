// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kubeadm // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm"

import kubeadmconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm/config"

type Config = kubeadmconfig.Config

func CreateDefaultConfig() Config {
	return kubeadmconfig.CreateDefaultConfig()
}
