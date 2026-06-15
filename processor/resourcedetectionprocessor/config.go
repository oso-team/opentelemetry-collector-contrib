// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

import (
	"time"

	"go.opentelemetry.io/collector/config/confighttp"

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

// Config defines configuration for Resource processor.
type Config struct {
	// Detectors is an ordered list of named detectors that should be
	// run to attempt to detect resource information.
	Detectors []string `mapstructure:"detectors"`
	// Override indicates whether any existing resource attributes
	// should be overridden or preserved. Defaults to true.
	Override bool `mapstructure:"override"`
	// DetectorConfig is a list of settings specific to all detectors
	DetectorConfig DetectorConfig `mapstructure:",squash"`
	// HTTP client settings for the detector
	// Timeout default is 5s
	confighttp.ClientConfig `mapstructure:",squash"`
	// If > 0, periodically re-run detection for all configured detectors.
	// When 0 (default), no periodic refresh occurs.
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`
}

// DetectorConfig contains user-specified configurations unique to all individual detectors
type DetectorConfig struct {
	// AlibabaECSConfig contains user-specified configurations for the Alibaba Cloud ECS detector
	AlibabaECSConfig alibabaecs.Config `mapstructure:"alibaba_ecs"`

	// EC2Config contains user-specified configurations for the EC2 detector
	EC2Config ec2.Config `mapstructure:"ec2"`

	// ECSConfig contains user-specified configurations for the ECS detector
	ECSConfig ecs.Config `mapstructure:"ecs"`

	// EKSConfig contains user-specified configurations for the EKS detector
	EKSConfig eks.Config `mapstructure:"eks"`

	// Elasticbeanstalk contains user-specified configurations for the elasticbeanstalk detector
	ElasticbeanstalkConfig elasticbeanstalk.Config `mapstructure:"elasticbeanstalk"`

	// Lambda contains user-specified configurations for the lambda detector
	LambdaConfig lambda.Config `mapstructure:"lambda"`

	// Azure contains user-specified configurations for the azure detector
	AzureConfig azure.Config `mapstructure:"azure"`

	// Aks contains user-specified configurations for the aks detector
	AksConfig aks.Config `mapstructure:"aks"`

	// ConsulConfig contains user-specified configurations for the Consul detector
	ConsulConfig consul.Config `mapstructure:"consul"`

	// DigitalOceanConfig contains user-specified configurations for the digitalocean detector
	DigitalOceanConfig digitalocean.Config `mapstructure:"digitalocean"`

	// DockerConfig contains user-specified configurations for the docker detector
	DockerConfig docker.Config `mapstructure:"docker"`

	// GcpConfig contains user-specified configurations for the gcp detector
	GcpConfig gcp.Config `mapstructure:"gcp"`

	// HerokuConfig contains user-specified configurations for the heroku detector
	HerokuConfig heroku.Config `mapstructure:"heroku"`

	// HetznerConfig contains user-specified configurations for the hetzner detector
	HetznerConfig hetzner.Config `mapstructure:"hetzner"`

	// IBMCloudClassicConfig contains user-specified configurations for the IBM Cloud Classic detector
	IBMCloudClassicConfig ibmcloudclassic.Config `mapstructure:"ibmcloud_classic"`

	// IBMCloudVPCConfig contains user-specified configurations for the IBM Cloud VPC detector
	IBMCloudVPCConfig ibmcloudvpc.Config `mapstructure:"ibmcloud_vpc"`

	// SystemConfig contains user-specified configurations for the System detector
	SystemConfig system.Config `mapstructure:"system"`

	// OpenShift contains user-specified configurations for the OpenShift detector
	OpenShiftConfig openshift.Config `mapstructure:"openshift"`

	// OpenStackNovaConfig contains user-specified configurations for the OpenStackNova detector
	OpenStackNovaConfig nova.Config `mapstructure:"nova"`

	// OracleCloud contains user-specified configurations for the OracleCloud detector
	OracleCloudConfig oraclecloud.Config `mapstructure:"oraclecloud"`

	// K8SAPIConfig contains user-specified configurations for the K8S API detector
	K8SAPIConfig k8sapi.Config `mapstructure:"k8s_api"`

	// K8SNodeConfig contains user-specified configurations for the K8SNode detector (deprecated: use K8SAPIConfig)
	K8SNodeConfig k8sapi.Config `mapstructure:"k8snode"`

	// Kubeadm contains user-specified configurations for the Kubeadm detector
	KubeadmConfig kubeadm.Config `mapstructure:"kubeadm"`

	// AkamaiConfig contains user-specified configurations for the akamai detector
	AkamaiConfig akamai.Config `mapstructure:"akamai"`

	// ScalewayConfig contains user-specified configurations for the scaleway detector
	ScalewayConfig scaleway.Config `mapstructure:"scaleway"`

	// TencentCVMConfig contains user-specified configurations for the Tencent Cloud CVM detector
	TencentCVMConfig tencentcvm.Config `mapstructure:"tencent_cvm"`

	// UpcloudConfig contains user-specified configurations for the upcloud detector
	UpcloudConfig upcloud.Config `mapstructure:"upcloud"`

	// VultrConfig contains user-specified configurations for the vultr detector
	VultrConfig vultr.Config `mapstructure:"vultr"`
}
