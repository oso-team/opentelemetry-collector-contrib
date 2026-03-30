package oraclecloudmonitoringexporter

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/oraclecloudauthextension"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

type oracleCloudMonitoringClient struct {
	logger *zap.Logger
	client monitoringClient
}

type monitoringClient interface {
	PostMetricData(ctx context.Context, request monitoring.PostMetricDataRequest) (response monitoring.PostMetricDataResponse, err error)
}

var (
	newMonitoringClient = func(provider common.ConfigurationProvider) (monitoring.MonitoringClient, error) {
		return monitoring.NewMonitoringClientWithConfigurationProvider(provider)
	}
)

func newMonitoringClientFromHost(cfg *Config, logger *zap.Logger, host component.Host) (*oracleCloudMonitoringClient, error) {
	provider, err := resolveAuthExtensionProvider(cfg, host)
	if err != nil {
		return nil, err
	}
	return newMonitoringClientWithProvider(cfg, logger, provider)
}

func newMonitoringClientWithProvider(cfg *Config, logger *zap.Logger, provider common.ConfigurationProvider) (*oracleCloudMonitoringClient, error) {
	monitoringClient, err := newMonitoringClient(provider)
	if err != nil {
		return nil, fmt.Errorf("failed creating monitoring client: %w", err)
	}

	monitoringClient.SetRegion(cfg.Region)

	if monitoringClient.Host != "" {
		monitoringClient.Host = strings.Replace(monitoringClient.Host, "telemetry.", "telemetry-ingestion.", 1)
	}

	return &oracleCloudMonitoringClient{
		logger: logger,
		client: &monitoringClient,
	}, nil
}

func resolveAuthExtensionProvider(cfg *Config, host component.Host) (common.ConfigurationProvider, error) {
	authID := cfg.Auth.Get().AuthenticatorID
	authExt, found := host.GetExtensions()[authID]
	if !found {
		return nil, fmt.Errorf("auth extension %q was not found", authID.String())
	}

	providerExt, ok := authExt.(oraclecloudauthextension.ConfigurationProvider)
	if !ok {
		return nil, fmt.Errorf("auth extension %q does not implement oraclecloudauth configuration provider", authID.String())
	}

	provider := providerExt.ConfigurationProvider()
	if provider == nil {
		return nil, fmt.Errorf("auth extension %q returned nil configuration provider", authID.String())
	}
	return provider, nil
}

func (c *oracleCloudMonitoringClient) SendMetrics(ctx context.Context, metricData []monitoring.MetricDataDetails) error {
	resp, err := c.client.PostMetricData(ctx, monitoring.PostMetricDataRequest{
		PostMetricDataDetails: monitoring.PostMetricDataDetails{
			MetricData: metricData,
		},
	})
	if err != nil {
		return fmt.Errorf("PostMetricData failed: %w", err)
	}

	// request got accepted with failures
	if resp.FailedMetricsCount != nil && *resp.FailedMetricsCount > 0 {
		// no handling. drop
	}

	return nil
}

func isPermanentMonitoringError(err error) bool {
	var serviceErr common.ServiceError
	if !errors.As(err, &serviceErr) {
		return false
	}

	status := serviceErr.GetHTTPStatusCode()
	if status >= 500 {
		return false
	}
	if status == 429 {
		return false
	}

	return status >= 400 && status < 500
}
