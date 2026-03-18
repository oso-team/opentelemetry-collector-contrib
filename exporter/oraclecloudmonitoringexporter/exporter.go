package oraclecloudmonitoringexporter

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type metricsExporter struct {
	cfg    *Config
	logger *zap.Logger
	client *ociClient
}

func (e *metricsExporter) start(_ context.Context, host component.Host) error {
	client, err := newOCIClientFromHost(e.cfg, e.logger, host)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

// pushMetricsData is called by exporterhelper to export metrics to OCI Monitoring.
func (e *metricsExporter) pushMetricsData(ctx context.Context, md pmetric.Metrics) error {
	if e.client == nil {
		return errors.New("OCI client was not initialized")
	}

	metricData, dropped, err := translateMetrics(md)
	if err != nil {
		return fmt.Errorf("failed to translate metrics: %w", err)
	}
	if len(metricData) == 0 {
		return nil
	}

	if err := e.client.SendMetrics(ctx, metricData); err != nil {
		return fmt.Errorf("failed to send metrics to OCI Monitoring: %w", err)
	}

	if dropped == 0 {
		return nil
	}

	e.logger.Debug("dropped unsupported metric datapoints during OCI translation", zap.Int("dropped_datapoints", dropped))
	return nil
}
