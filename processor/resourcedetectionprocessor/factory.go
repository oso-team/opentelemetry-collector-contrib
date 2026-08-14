// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resourcedetectionprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.opentelemetry.io/collector/processor/processorhelper/xprocessorhelper"
	"go.opentelemetry.io/collector/processor/xprocessor"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/metadata"
)

var consumerCapabilities = consumer.Capabilities{MutatesData: true}

type factory struct {
	resourceProviderFactory *internal.ResourceProviderFactory
	compiledDetectors       []string

	// providers stores a provider for each named processor that
	// may a different set of detectors configured.
	providers map[component.ID]*internal.ResourceProvider
	lock      sync.Mutex
}

// NewFactory creates a new factory for ResourceDetection processor.
func NewFactory() processor.Factory {
	f := &factory{
		resourceProviderFactory: internal.NewProviderFactory(globalDetectorRegistry),
		compiledDetectors:       sortedDetectorTypes(globalDetectorRegistry),
		providers:               map[component.ID]*internal.ResourceProvider{},
	}

	return xprocessor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		xprocessor.WithDeprecatedTypeAlias(metadata.DeprecatedType),
		xprocessor.WithTraces(f.createTracesProcessor, metadata.TracesStability),
		xprocessor.WithMetrics(f.createMetricsProcessor, metadata.MetricsStability),
		xprocessor.WithLogs(f.createLogsProcessor, metadata.LogsStability),
		xprocessor.WithProfiles(f.createProfilesProcessor, metadata.ProfilesStability),
	)
}

// Type gets the type of the Option config created by this factory.
func (*factory) Type() component.Type {
	return metadata.Type
}

func createDefaultConfig() component.Config {
	return &Config{
		Detectors:       []string{"env"},
		ClientConfig:    defaultClientConfig(),
		Override:        true,
		DetectorConfig:  detectorCreateDefaultConfig(),
		RefreshInterval: 0,
		Retry:           defaultRetryConfig(),
		// TODO: Once issue(https://github.com/open-telemetry/opentelemetry-collector/issues/4001) gets resolved,
		//		 Set the default value of 'hostname_source' here instead of 'system' detector
	}
}

func defaultRetryConfig() configretry.BackOffConfig {
	return configretry.BackOffConfig{
		Enabled:             true,
		InitialInterval:     1 * time.Second,
		RandomizationFactor: 0.5,
		Multiplier:          2,
		MaxInterval:         30 * time.Second,
		MaxElapsedTime:      0,
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

	warnDeprecatedPerDetectorFlags(params.Logger, oCfg)

	// The deprecated per-detector fail_on_missing_metadata fields are OR'd with this
	// top-level flag inside each affected detector, preserving per-detector scope.
	provider, err := f.getResourceProvider(params, oCfg.Retry, oCfg.Detectors, oCfg.DetectorConfig, oCfg.FailOnMissingMetadata)
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

// warnDeprecatedPerDetectorFlags emits a deprecation warning if any per-detector
// fail_on_missing_metadata fields are set.
func warnDeprecatedPerDetectorFlags(logger *zap.Logger, oCfg *Config) {
	var affectedDetectors []string
	if oCfg.DetectorConfig.EC2Config.FailOnMissingMetadata { //nolint:staticcheck
		affectedDetectors = append(affectedDetectors, "ec2")
	}
	if oCfg.DetectorConfig.AlibabaECSConfig.FailOnMissingMetadata { //nolint:staticcheck
		affectedDetectors = append(affectedDetectors, "alibaba_ecs")
	}
	if oCfg.DetectorConfig.TencentCVMConfig.FailOnMissingMetadata { //nolint:staticcheck
		affectedDetectors = append(affectedDetectors, "tencent_cvm")
	}
	if oCfg.DetectorConfig.UpcloudConfig.FailOnMissingMetadata { //nolint:staticcheck
		affectedDetectors = append(affectedDetectors, "upcloud")
	}
	if oCfg.DetectorConfig.VultrConfig.FailOnMissingMetadata { //nolint:staticcheck
		affectedDetectors = append(affectedDetectors, "vultr")
	}
	if oCfg.DetectorConfig.OpenStackNovaConfig.FailOnMissingMetadata { //nolint:staticcheck
		affectedDetectors = append(affectedDetectors, "nova")
	}
	if len(affectedDetectors) > 0 {
		logger.Warn(
			"per-detector fail_on_missing_metadata fields are deprecated; use the top-level fail_on_missing_metadata instead",
			zap.Strings("affected_detectors", affectedDetectors),
		)
	}
}

func (f *factory) getResourceProvider(
	params processor.Settings,
	backoffConfig configretry.BackOffConfig,
	configuredDetectors []string,
	detectorConfigs DetectorConfig,
	failOnMissingMetadata bool,
) (*internal.ResourceProvider, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if provider, ok := f.providers[params.ID]; ok {
		return provider, nil
	}

	detectorTypes := make([]internal.DetectorType, 0, len(configuredDetectors))
	configuredDetectorTypes := make([]string, 0, len(configuredDetectors))
	for _, key := range configuredDetectors {
		detectorType := strings.TrimSpace(key)
		detectorTypes = append(detectorTypes, internal.DetectorType(detectorType))
		configuredDetectorTypes = append(configuredDetectorTypes, detectorType)
	}

	logDetectorConfiguration(params.Logger, f.compiledDetectors, configuredDetectorTypes)

	provider, err := f.resourceProviderFactory.CreateResourceProvider(params, backoffConfig, failOnMissingMetadata, &detectorConfigs, detectorTypes...)
	if err != nil {
		return nil, err
	}

	f.providers[params.ID] = provider
	return provider, nil
}

func sortedDetectorTypes(registry map[internal.DetectorType]internal.DetectorFactory) []string {
	detectorTypes := make([]string, 0, len(registry))
	for detectorType := range registry {
		detectorTypes = append(detectorTypes, string(detectorType))
	}

	sort.Strings(detectorTypes)
	return detectorTypes
}
