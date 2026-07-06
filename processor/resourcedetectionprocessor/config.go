// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

import (
	"time"

	"go.opentelemetry.io/collector/config/confighttp"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
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

// DetectorConfig contains user-specified configurations unique to all individual detectors.
// The field types come from each detector's dependency-light config package rather than the
// detector package itself, so that this always-compiled schema does not import any detector
// implementation (and the vendor SDKs they carry).
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

func detectorCreateDefaultConfig() DetectorConfig {
	return DetectorConfig{
		AlibabaECSConfig:       alibabaecsconfig.CreateDefaultConfig(),
		EC2Config:              ec2config.CreateDefaultConfig(),
		ECSConfig:              ecsconfig.CreateDefaultConfig(),
		EKSConfig:              eksconfig.CreateDefaultConfig(),
		ElasticbeanstalkConfig: elasticbeanstalkconfig.CreateDefaultConfig(),
		LambdaConfig:           lambdaconfig.CreateDefaultConfig(),
		AzureConfig:            azureconfig.CreateDefaultConfig(),
		AksConfig:              aksconfig.CreateDefaultConfig(),
		ConsulConfig:           consulconfig.CreateDefaultConfig(),
		DigitalOceanConfig:     digitaloceanconfig.CreateDefaultConfig(),
		DockerConfig:           dockerconfig.CreateDefaultConfig(),
		GcpConfig:              gcpconfig.CreateDefaultConfig(),
		HerokuConfig:           herokuconfig.CreateDefaultConfig(),
		HetznerConfig:          hetznerconfig.CreateDefaultConfig(),
		IBMCloudClassicConfig:  ibmcloudclassicconfig.CreateDefaultConfig(),
		IBMCloudVPCConfig:      ibmcloudvpcconfig.CreateDefaultConfig(),
		SystemConfig:           systemconfig.CreateDefaultConfig(),
		OpenShiftConfig:        openshiftconfig.CreateDefaultConfig(),
		OpenStackNovaConfig:    novaconfig.CreateDefaultConfig(),
		OracleCloudConfig:      oraclecloudconfig.CreateDefaultConfig(),
		K8SAPIConfig:           k8sapiconfig.CreateDefaultConfig(),
		K8SNodeConfig:          k8sapiconfig.CreateDefaultConfig(),
		KubeadmConfig:          kubeadmconfig.CreateDefaultConfig(),
		AkamaiConfig:           akamaiconfig.CreateDefaultConfig(),
		ScalewayConfig:         scalewayconfig.CreateDefaultConfig(),
		TencentCVMConfig:       tencentcvmconfig.CreateDefaultConfig(),
		UpcloudConfig:          upcloudconfig.CreateDefaultConfig(),
		VultrConfig:            vultrconfig.CreateDefaultConfig(),
	}
}

// GetConfigFromType returns the per-detector config for the given detector type.
// The case values are the detectors' type strings (the same strings used as
// mapstructure tags above and by the register_*.go files); string literals are
// used instead of the detector packages' TypeStr constants so that this
// always-compiled file does not import any detector implementation.
func (d *DetectorConfig) GetConfigFromType(detectorType internal.DetectorType) internal.DetectorConfig {
	switch detectorType {
	case "alibaba_ecs":
		return d.AlibabaECSConfig
	case "ec2":
		return d.EC2Config
	case "ecs":
		return d.ECSConfig
	case "eks":
		return d.EKSConfig
	case "elastic_beanstalk":
		return d.ElasticbeanstalkConfig
	case "lambda":
		return d.LambdaConfig
	case "azure":
		return d.AzureConfig
	case "aks":
		return d.AksConfig
	case "consul":
		return d.ConsulConfig
	case "digitalocean":
		return d.DigitalOceanConfig
	case "docker":
		return d.DockerConfig
	case "gcp":
		return d.GcpConfig
	case "heroku":
		return d.HerokuConfig
	case "hetzner":
		return d.HetznerConfig
	case "ibmcloud_classic":
		return d.IBMCloudClassicConfig
	case "ibmcloud_vpc":
		return d.IBMCloudVPCConfig
	case "system":
		return d.SystemConfig
	case "openshift":
		return d.OpenShiftConfig
	case "nova":
		return d.OpenStackNovaConfig
	case "oraclecloud":
		return d.OracleCloudConfig
	case "k8s_api":
		return d.K8SAPIConfig
	case "k8snode":
		return d.K8SNodeConfig
	case "kubeadm":
		return d.KubeadmConfig
	case "akamai":
		return d.AkamaiConfig
	case "scaleway":
		return d.ScalewayConfig
	case "tencent_cvm":
		return d.TencentCVMConfig
	case "upcloud":
		return d.UpcloudConfig
	case "vultr":
		return d.VultrConfig
	default:
		return nil
	}
}
