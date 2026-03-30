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

	metricData, translationDropped, err := translateMetrics(md, e.cfg.CompartmentId, e.cfg.Namespace)
	if translationDropped > 0 {
		// e.logger.Debug("dropped metric datapoints during metric translation", zap.Int("dropped_datapoints", translationDropped))
	}
	if err != nil {
		return consumererror.NewPermanent(fmt.Errorf("failed to translate metrics: %w", err))
	}
	// batch metrics first to meet per call limit of metrics stream
	batchedMetrics := buildMetricBatches(metricData)
	if len(batchedMetrics.batches) == 0 {
		return nil
	}

	if len(batchedMetrics.batches) > 1 {
		// e.logger.Debug("split metrics into multiple ingestion requests", zap.Int("batched_requests", len(batchedMetrics.batches)))
	}

	// Send batches sequentially
	for i, batch := range batchedMetrics.batches {
		if err := e.client.SendMetrics(ctx, batch); err != nil {
			sendErr := fmt.Errorf("failed to send metrics to Oracle Cloud Monitoring: %w", err)
			return consumererror.NewPermanent(sendErr)
		}

		if len(batchedMetrics.batches) > 1 {
			e.logger.Debug("sent monitoring request", zap.Int("request_index", i+1))
		}
	}

	return nil
}
