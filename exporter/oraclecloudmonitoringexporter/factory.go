//go:generate make mdatagen

package oraclecloudmonitoringexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

var (
	Type = component.MustNewType("oraclecloudmonitoring")
)

// NewFactory creates a factory for the Oracle Cloud Monitoring exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		Type,
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, component.StabilityLevelAlpha),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		MaxPastTimestampSkew:   defaultMaxPastTimestampSkew,
		MaxFutureTimestampSkew: defaultMaxFutureTimestampSkew,
		TimeoutSettings:        exporterhelper.NewDefaultTimeoutConfig(),
		RetrySettings:          configretry.NewDefaultBackOffConfig(),
		QueueSettings:          configoptional.Some(exporterhelper.NewDefaultQueueConfig()),
	}
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Metrics, error) {
	exporterCfg := cfg.(*Config)

	if err := exporterCfg.Validate(); err != nil {
		return nil, err
	}

	exp := &metricsExporter{
		cfg:    exporterCfg,
		logger: set.Logger,
	}

	return exporterhelper.NewMetrics(
		ctx,
		set,
		cfg,
		exp.pushMetricsData,
		exporterhelper.WithStart(exp.start),
		exporterhelper.WithTimeout(exporterCfg.TimeoutSettings),
		exporterhelper.WithRetry(exporterCfg.RetrySettings),
		exporterhelper.WithQueue(exporterCfg.QueueSettings),
	)
}
