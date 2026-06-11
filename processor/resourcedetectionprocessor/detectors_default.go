// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/akamai"
	alibabaecs "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/alibaba/ecs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ec2"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ecs"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/eks"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/elasticbeanstalk"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/lambda"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/aks"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/consul"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/digitalocean"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/docker"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/dynatrace"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/env"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/gcp"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/heroku"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/hetzner"
	ibmcloudclassic "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/ibmcloud/classic"
	ibmcloudvpc "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/ibmcloud/vpc"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/k8sapi"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openshift"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openstack/nova"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/oraclecloud"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/scaleway"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/system"
	tencentcvm "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/tencent/cvm"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/upcloud"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/vultr"
)

func init() {
	registerDetector(akamai.TypeStr, akamai.NewDetector)
	registerDetector(alibabaecs.TypeStr, alibabaecs.NewDetector)
	registerDetector(aks.TypeStr, aks.NewDetector)
	registerDetector(azure.TypeStr, azure.NewDetector)
	registerDetector(consul.TypeStr, consul.NewDetector)
	registerDetector(digitalocean.TypeStr, digitalocean.NewDetector)
	registerDetector(docker.TypeStr, docker.NewDetector)
	registerDetector(ec2.TypeStr, ec2.NewDetector)
	registerDetector(ecs.TypeStr, ecs.NewDetector)
	registerDetector(eks.TypeStr, eks.NewDetector)
	registerDetector(elasticbeanstalk.TypeStr, elasticbeanstalk.NewDetector)
	registerDetector(lambda.TypeStr, lambda.NewDetector)
	registerDetector(env.TypeStr, env.NewDetector)
	registerDetector(gcp.TypeStr, gcp.NewDetector)
	registerDetector(heroku.TypeStr, heroku.NewDetector)
	registerDetector(hetzner.TypeStr, hetzner.NewDetector)
	registerDetector(ibmcloudclassic.TypeStr, ibmcloudclassic.NewDetector)
	registerDetector(ibmcloudvpc.TypeStr, ibmcloudvpc.NewDetector)
	registerDetector(scaleway.TypeStr, scaleway.NewDetector)
	registerDetector(system.TypeStr, system.NewDetector)
	registerDetector(openshift.TypeStr, openshift.NewDetector)
	registerDetector(nova.TypeStr, nova.NewDetector)
	registerDetector(oraclecloud.TypeStr, oraclecloud.NewDetector)
	registerDetector(k8sapi.TypeStr, k8sapi.NewDetector)
	registerDetector(k8sapi.TypeStrAlias, k8sapi.NewDeprecatedDetector)
	registerDetector(kubeadm.TypeStr, kubeadm.NewDetector)
	registerDetector(dynatrace.TypeStr, dynatrace.NewDetector)
	registerDetector(tencentcvm.TypeStr, tencentcvm.NewDetector)
	registerDetector(upcloud.TypeStr, upcloud.NewDetector)
	registerDetector(vultr.TypeStr, vultr.NewDetector)
}
