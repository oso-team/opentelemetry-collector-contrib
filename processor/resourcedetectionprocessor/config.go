// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

import (
	"time"

	"go.opentelemetry.io/collector/config/confighttp"

	akamaiconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/akamai/config"
	alibabaecsconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/alibaba/ecs/config"
	ec2config "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ec2/config"
	ecsconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/ecs/config"
	eksconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/eks/config"
	elasticbeanstalkconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/elasticbeanstalk/config"
	lambdaconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/aws/lambda/config"
	aksconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/aks/config"
	azureconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/azure/config"
	consulconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/consul/config"
	digitaloceanconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/digitalocean/config"
	dockerconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/docker/config"
	gcpconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/gcp/config"
	herokuconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/heroku/config"
	hetznerconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/hetzner/config"
	ibmcloudclassicconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/ibmcloud/classic/config"
	ibmcloudvpcconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/ibmcloud/vpc/config"
	k8sapiconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/k8sapi/config"
	kubeadmconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/kubeadm/config"
	openshiftconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openshift/config"
	novaconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openstack/nova/config"
	oraclecloudconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/oraclecloud/config"
	scalewayconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/scaleway/config"
	systemconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/system/config"
	tencentcvmconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/tencent/cvm/config"
	upcloudconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/upcloud/config"
	vultrconfig "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/vultr/config"
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
	AlibabaECSConfig alibabaecsconfig.Config `mapstructure:"alibaba_ecs"`

	// EC2Config contains user-specified configurations for the EC2 detector
	EC2Config ec2config.Config `mapstructure:"ec2"`

	// ECSConfig contains user-specified configurations for the ECS detector
	ECSConfig ecsconfig.Config `mapstructure:"ecs"`

	// EKSConfig contains user-specified configurations for the EKS detector
	EKSConfig eksconfig.Config `mapstructure:"eks"`

	// Elasticbeanstalk contains user-specified configurations for the elasticbeanstalk detector
	ElasticbeanstalkConfig elasticbeanstalkconfig.Config `mapstructure:"elasticbeanstalk"`

	// Lambda contains user-specified configurations for the lambda detector
	LambdaConfig lambdaconfig.Config `mapstructure:"lambda"`

	// Azure contains user-specified configurations for the azure detector
	AzureConfig azureconfig.Config `mapstructure:"azure"`

	// Aks contains user-specified configurations for the aks detector
	AksConfig aksconfig.Config `mapstructure:"aks"`

	// ConsulConfig contains user-specified configurations for the Consul detector
	ConsulConfig consulconfig.Config `mapstructure:"consul"`

	// DigitalOceanConfig contains user-specified configurations for the digitalocean detector
	DigitalOceanConfig digitaloceanconfig.Config `mapstructure:"digitalocean"`

	// DockerConfig contains user-specified configurations for the docker detector
	DockerConfig dockerconfig.Config `mapstructure:"docker"`

	// GcpConfig contains user-specified configurations for the gcp detector
	GcpConfig gcpconfig.Config `mapstructure:"gcp"`

	// HerokuConfig contains user-specified configurations for the heroku detector
	HerokuConfig herokuconfig.Config `mapstructure:"heroku"`

	// HetznerConfig contains user-specified configurations for the hetzner detector
	HetznerConfig hetznerconfig.Config `mapstructure:"hetzner"`

	// IBMCloudClassicConfig contains user-specified configurations for the IBM Cloud Classic detector
	IBMCloudClassicConfig ibmcloudclassicconfig.Config `mapstructure:"ibmcloud_classic"`

	// IBMCloudVPCConfig contains user-specified configurations for the IBM Cloud VPC detector
	IBMCloudVPCConfig ibmcloudvpcconfig.Config `mapstructure:"ibmcloud_vpc"`

	// SystemConfig contains user-specified configurations for the System detector
	SystemConfig systemconfig.Config `mapstructure:"system"`

	// OpenShift contains user-specified configurations for the OpenShift detector
	OpenShiftConfig openshiftconfig.Config `mapstructure:"openshift"`

	// OpenStackNovaConfig contains user-specified configurations for the OpenStackNova detector
	OpenStackNovaConfig novaconfig.Config `mapstructure:"nova"`

	// OracleCloud contains user-specified configurations for the OracleCloud detector
	OracleCloudConfig oraclecloudconfig.Config `mapstructure:"oraclecloud"`

	// K8SAPIConfig contains user-specified configurations for the K8S API detector
	K8SAPIConfig k8sapiconfig.Config `mapstructure:"k8s_api"`

	// K8SNodeConfig contains user-specified configurations for the K8SNode detector (deprecated: use K8SAPIConfig)
	K8SNodeConfig k8sapiconfig.Config `mapstructure:"k8snode"`

	// Kubeadm contains user-specified configurations for the Kubeadm detector
	KubeadmConfig kubeadmconfig.Config `mapstructure:"kubeadm"`

	// AkamaiConfig contains user-specified configurations for the akamai detector
	AkamaiConfig akamaiconfig.Config `mapstructure:"akamai"`

	// ScalewayConfig contains user-specified configurations for the scaleway detector
	ScalewayConfig scalewayconfig.Config `mapstructure:"scaleway"`

	// TencentCVMConfig contains user-specified configurations for the Tencent Cloud CVM detector
	TencentCVMConfig tencentcvmconfig.Config `mapstructure:"tencent_cvm"`

	// UpcloudConfig contains user-specified configurations for the upcloud detector
	UpcloudConfig upcloudconfig.Config `mapstructure:"upcloud"`

	// VultrConfig contains user-specified configurations for the vultr detector
	VultrConfig vultrconfig.Config `mapstructure:"vultr"`
}
