package oraclecloudmonitoringexporter

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type metricsExporter struct {
	cfg    *Config
	logger *zap.Logger
	client *oracleCloudMonitoringClient
}

func (e *metricsExporter) start(_ context.Context, host component.Host) error {
	client, err := newMonitoringClientFromHost(e.cfg, e.logger, host)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

// pushMetricsData is called by exporterhelper to export metrics to Oracle Cloud Monitoring.
func (e *metricsExporter) pushMetricsData(ctx context.Context, md pmetric.Metrics) error {
	if e.client == nil {
		return consumererror.NewPermanent(errors.New("Monitoring client was not initialized"))
	}

	metricData, dropped, err := translateMetrics(md, e.cfg.CompartmentId, e.cfg.Namespace)
	if err != nil {
		return consumererror.NewPermanent(fmt.Errorf("failed to translate metrics: %w", err))
	}
	if len(metricData) == 0 {
		return nil
	}

	if err := e.client.SendMetrics(ctx, metricData); err != nil {
		sendErr := fmt.Errorf("failed to send metrics to Oracle Cloud Monitoring: %w", err)
		if isPermanentMonitoringError(err) {
			return consumererror.NewPermanent(sendErr)
		}
		return sendErr
	}

	if dropped == 0 {
		return nil
	}

	e.logger.Debug("dropped unsupported metric datapoints during metric translation", zap.Int("dropped_datapoints", dropped))
	return nil
}
