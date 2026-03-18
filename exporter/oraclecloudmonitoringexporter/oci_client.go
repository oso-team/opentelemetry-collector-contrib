package oraclecloudmonitoringexporter

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/oraclecloudauthextension"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

type ociClient struct {
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

func newOCIClientFromHost(cfg *Config, logger *zap.Logger, host component.Host) (*ociClient, error) {
	provider, err := resolveAuthExtensionProvider(cfg, host)
	if err != nil {
		return nil, err
	}
	return newOCIClientWithProvider(cfg, logger, provider)
}

func newOCIClientWithProvider(cfg *Config, logger *zap.Logger, provider common.ConfigurationProvider) (*ociClient, error) {
	ociMonitoringClient, err := newMonitoringClient(provider)
	if err != nil {
		return nil, fmt.Errorf("failed creating OCI Monitoring client: %w", err)
	}

	ociMonitoringClient.SetRegion(cfg.Region)

	if ociMonitoringClient.Host != "" {
		ociMonitoringClient.Host = strings.Replace(ociMonitoringClient.Host, "telemetry.", "telemetry-ingestion.", 1)
	}

	return &ociClient{
		logger: logger,
		client: &ociMonitoringClient,
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
		return nil, fmt.Errorf("auth extension %q returned nil OCI configuration provider", authID.String())
	}
	return provider, nil
}

func (c *ociClient) SendMetrics(ctx context.Context, metricData []monitoring.MetricDataDetails) error {
	_, err := c.client.PostMetricData(ctx, monitoring.PostMetricDataRequest{
		PostMetricDataDetails: monitoring.PostMetricDataDetails{
			MetricData: metricData,
		},
	})
	if err != nil {
		return fmt.Errorf("PostMetricData failed: %w", err)
	}
	return nil
}
