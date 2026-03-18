package oraclecloudmonitoringexporter

import (
	"context"
	"errors"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/oraclecloudauthextension"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.uber.org/zap"
)

type fakeMonitoringClient struct {
	lastRequest monitoring.PostMetricDataRequest
	err         error
}

func (f *fakeMonitoringClient) PostMetricData(_ context.Context, request monitoring.PostMetricDataRequest) (monitoring.PostMetricDataResponse, error) {
	f.lastRequest = request
	if f.err != nil {
		return monitoring.PostMetricDataResponse{}, f.err
	}
	return monitoring.PostMetricDataResponse{}, nil
}

type fakeHost struct {
	extensions map[component.ID]component.Component
}

func (f fakeHost) GetExtensions() map[component.ID]component.Component {
	return f.extensions
}

func (fakeHost) GetFactory(component.Kind, component.Type) component.Factory {
	return nil
}

type fakeOracleCloudAuthExtension struct {
	component.StartFunc
	component.ShutdownFunc
	provider common.ConfigurationProvider
}

func (f *fakeOracleCloudAuthExtension) ConfigurationProvider() common.ConfigurationProvider {
	return f.provider
}

var _ oraclecloudauthextension.ConfigurationProvider = (*fakeOracleCloudAuthExtension)(nil)

func TestSendMetrics(t *testing.T) {
	fake := &fakeMonitoringClient{}
	client := &ociClient{
		logger: zap.NewNop(),
		client: fake,
	}

	name := "cpu.utilization"
	namespace := "otel_demo"
	compartment := "ocid1.compartment.oc1..aaaa"
	err := client.SendMetrics(t.Context(), []monitoring.MetricDataDetails{{
		Name:          &name,
		Namespace:     &namespace,
		CompartmentId: &compartment,
	}})
	require.NoError(t, err)
	require.Len(t, fake.lastRequest.PostMetricDataDetails.MetricData, 1)
}

func TestSendMetricsError(t *testing.T) {
	fake := &fakeMonitoringClient{err: errors.New("boom")}
	client := &ociClient{
		logger: zap.NewNop(),
		client: fake,
	}

	err := client.SendMetrics(t.Context(), []monitoring.MetricDataDetails{{}})
	require.Error(t, err)
}

func TestResolveAuthExtensionProvider(t *testing.T) {
	authID := component.MustNewID("oraclecloudauth")
	cfg := &Config{
		Region: "us-phoenix-1",
		Auth: configoptional.Some(configauth.Config{
			AuthenticatorID: authID,
		}),
	}

	provider, err := resolveAuthExtensionProvider(cfg, fakeHost{
		extensions: map[component.ID]component.Component{
			authID: &fakeOracleCloudAuthExtension{provider: common.DefaultConfigProvider()},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, provider)
}

func TestResolveAuthExtensionProviderMissing(t *testing.T) {
	authID := component.MustNewID("oraclecloudauth")
	cfg := &Config{
		Region: "us-phoenix-1",
		Auth: configoptional.Some(configauth.Config{
			AuthenticatorID: authID,
		}),
	}

	_, err := resolveAuthExtensionProvider(cfg, fakeHost{extensions: map[component.ID]component.Component{}})
	require.Error(t, err)
}

func TestNewOCIClientFromHost(t *testing.T) {
	authID := component.MustNewID("oraclecloudauth")
	cfg := &Config{
		Region: "us-phoenix-1",
		Auth: configoptional.Some(configauth.Config{
			AuthenticatorID: authID,
		}),
	}

	client, err := newOCIClientFromHost(cfg, zap.NewNop(), fakeHost{
		extensions: map[component.ID]component.Component{
			authID: &fakeOracleCloudAuthExtension{provider: common.DefaultConfigProvider()},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}
