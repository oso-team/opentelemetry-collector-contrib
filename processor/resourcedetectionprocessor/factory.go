// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

import (
	"context"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.opentelemetry.io/collector/processor/processorhelper/xprocessorhelper"
	"go.opentelemetry.io/collector/processor/xprocessor"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
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
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openshift"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/openstack/nova"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/oraclecloud"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/scaleway"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/system"
	tencentcvm "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/tencent/cvm"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/upcloud"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/vultr"
)

var consumerCapabilities = consumer.Capabilities{MutatesData: true}

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
type factory struct {
	resourceProviderFactory *internal.ResourceProviderFactory

	// providers stores a provider for each named processor that
	// may a different set of detectors configured.
	providers map[component.ID]*internal.ResourceProvider
	lock      sync.Mutex
}

// NewFactory creates a new factory for ResourceDetection processor.
func NewFactory() processor.Factory {

	resourceProviderFactory := internal.NewProviderFactory(detectorRegistry())

	f := &factory{
		resourceProviderFactory: resourceProviderFactory,
		providers:               map[component.ID]*internal.ResourceProvider{},
	}

	return xprocessor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		xprocessor.WithDeprecatedTypeAlias(metadata.DeprecatedType),
		xprocessor.WithTraces(f.createTracesProcessor, metadata.TracesStability),
		xprocessor.WithMetrics(f.createMetricsProcessor, metadata.MetricsStability),
		xprocessor.WithLogs(f.createLogsProcessor, metadata.LogsStability),
		xprocessor.WithProfiles(f.createProfilesProcessor, metadata.ProfilesStability))
}

// Type gets the type of the Option config created by this factory.
func (*factory) Type() component.Type {
	return metadata.Type
}

func createDefaultConfig() component.Config {
	return &Config{
		Detectors:       []string{env.TypeStr},
		ClientConfig:    defaultClientConfig(),
		Override:        true,
		DetectorConfig:  detectorCreateDefaultConfig(),
		RefreshInterval: 0,
		// TODO: Once issue(https://github.com/open-telemetry/opentelemetry-collector/issues/4001) gets resolved,
		//		 Set the default value of 'hostname_source' here instead of 'system' detector
	}
}

func defaultClientConfig() confighttp.ClientConfig {
	httpClientSettings := confighttp.NewDefaultClientConfig()
	httpClientSettings.Timeout = 5 * time.Second
	return httpClientSettings
}

func (f *factory) createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	rdp, err := f.getResourceDetectionProcessor(set, cfg)
	if err != nil {
		return nil, err
	}

	return processorhelper.NewTraces(
		ctx,
		set,
		cfg,
		nextConsumer,
		rdp.processTraces,
		processorhelper.WithCapabilities(consumerCapabilities),
		processorhelper.WithStart(rdp.Start),
		processorhelper.WithShutdown(rdp.Shutdown),
	)
}

func (f *factory) createMetricsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (processor.Metrics, error) {
	rdp, err := f.getResourceDetectionProcessor(set, cfg)
	if err != nil {
		return nil, err
	}

	return processorhelper.NewMetrics(
		ctx,
		set,
		cfg,
		nextConsumer,
		rdp.processMetrics,
		processorhelper.WithCapabilities(consumerCapabilities),
		processorhelper.WithStart(rdp.Start),
		processorhelper.WithShutdown(rdp.Shutdown),
	)
}

func (f *factory) createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (processor.Logs, error) {
	rdp, err := f.getResourceDetectionProcessor(set, cfg)
	if err != nil {
		return nil, err
	}

	return processorhelper.NewLogs(
		ctx,
		set,
		cfg,
		nextConsumer,
		rdp.processLogs,
		processorhelper.WithCapabilities(consumerCapabilities),
		processorhelper.WithStart(rdp.Start),
		processorhelper.WithShutdown(rdp.Shutdown),
	)
}

func (f *factory) createProfilesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer xconsumer.Profiles,
) (xprocessor.Profiles, error) {
	rdp, err := f.getResourceDetectionProcessor(set, cfg)
	if err != nil {
		return nil, err
	}

	return xprocessorhelper.NewProfiles(
		ctx,
		set,
		cfg,
		nextConsumer,
		rdp.processProfiles,
		xprocessorhelper.WithCapabilities(consumerCapabilities),
		xprocessorhelper.WithStart(rdp.Start),
		xprocessorhelper.WithShutdown(rdp.Shutdown),
	)
}

func (f *factory) getResourceDetectionProcessor(
	params processor.Settings,
	cfg component.Config,
) (*resourceDetectionProcessor, error) {
	oCfg := cfg.(*Config)
	provider, err := f.getResourceProvider(params, oCfg.Timeout, oCfg.Detectors, oCfg.DetectorConfig)
	if err != nil {
		return nil, err
	}

	return &resourceDetectionProcessor{
		provider:           provider,
		override:           oCfg.Override,
		httpClientSettings: oCfg.ClientConfig,
		refreshInterval:    oCfg.RefreshInterval,
		telemetrySettings:  params.TelemetrySettings,
	}, nil
}

func (f *factory) getResourceProvider(
	params processor.Settings,
	timeout time.Duration,
	configuredDetectors []string,
	detectorConfigs DetectorConfig,
) (*internal.ResourceProvider, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if provider, ok := f.providers[params.ID]; ok {
		return provider, nil
	}

	detectorTypes := make([]internal.DetectorType, 0, len(configuredDetectors))
	for _, key := range configuredDetectors {
		detectorTypes = append(detectorTypes, internal.DetectorType(strings.TrimSpace(key)))
	}

	provider, err := f.resourceProviderFactory.CreateResourceProvider(params, timeout, &detectorConfigs, detectorTypes...)
	if err != nil {
		return nil, err
	}

	f.providers[params.ID] = provider
	return provider, nil
}
