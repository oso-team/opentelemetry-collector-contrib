// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
)

func TestDefaultDetectorRegistry(t *testing.T) {
	expectedTypes := []internal.DetectorType{
		"akamai",
		"alibaba_ecs",
		"aks",
		"azure",
		"consul",
		"digitalocean",
		"docker",
		"dynatrace",
		"ec2",
		"ecs",
		"eks",
		"elastic_beanstalk",
		"env",
		"gcp",
		"heroku",
		"hetzner",
		"ibmcloud_classic",
		"ibmcloud_vpc",
		"k8s_api",
		"k8snode",
		"kubeadm",
		"lambda",
		"nova",
		"openshift",
		"oraclecloud",
		"scaleway",
		"system",
		"tencent_cvm",
		"upcloud",
		"vultr",
	}

	registry := detectorRegistry()
	actualTypes := make([]internal.DetectorType, 0, len(registry))
	for detectorType, factory := range registry {
		require.NotNil(t, factory, "factory should not be nil for %q", detectorType)
		actualTypes = append(actualTypes, detectorType)
	}

	assert.ElementsMatch(t, expectedTypes, actualTypes)
}

func TestDetectorRegistryReturnsCopy(t *testing.T) {
	registry := detectorRegistry()
	delete(registry, "env")

	assert.Contains(t, globalDetectorRegistry, internal.DetectorType("env"))
}
