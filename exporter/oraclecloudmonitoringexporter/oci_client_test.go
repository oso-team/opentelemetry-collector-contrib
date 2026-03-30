package oraclecloudmonitoringexporter

import (
	"context"
	"errors"
	"fmt"
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
	resp        monitoring.PostMetricDataResponse
	err         error
}

func (f *fakeMonitoringClient) PostMetricData(_ context.Context, request monitoring.PostMetricDataRequest) (monitoring.PostMetricDataResponse, error) {
	f.lastRequest = request
	if f.err != nil {
		return monitoring.PostMetricDataResponse{}, f.err
	}
	return f.resp, nil
}

type fakeServiceError struct {
	statusCode int
	code       string
	message    string
	opcReqID   string
}

func (e fakeServiceError) Error() string {
	return e.message
}

func (e fakeServiceError) GetHTTPStatusCode() int {
	return e.statusCode
}

func (e fakeServiceError) GetMessage() string {
	return e.message
}

func (e fakeServiceError) GetCode() string {
	return e.code
}

func (e fakeServiceError) GetOpcRequestID() string {
	return e.opcReqID
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
	client := &oracleCloudMonitoringClient{
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

func TestSendMetricsPartialValidationFailure(t *testing.T) {
	failed := 2
	fake := &fakeMonitoringClient{
		resp: monitoring.PostMetricDataResponse{
			PostMetricDataResponseDetails: monitoring.PostMetricDataResponseDetails{
				FailedMetricsCount: &failed,
			},
		},
	}
	client := &oracleCloudMonitoringClient{
		logger: zap.NewNop(),
		client: fake,
	}

	err := client.SendMetrics(t.Context(), []monitoring.MetricDataDetails{{}})
	require.NoError(t, err)
}

func TestSendMetricsError(t *testing.T) {
	fake := &fakeMonitoringClient{err: errors.New("boom")}
	client := &oracleCloudMonitoringClient{
		logger: zap.NewNop(),
		client: fake,
	}

	err := client.SendMetrics(t.Context(), []monitoring.MetricDataDetails{{}})
	require.Error(t, err)
}

func TestIsPermanentMonitoringError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		permanent bool
	}{
		{
			name: "bad request 400 is permanent",
			err: fakeServiceError{
				statusCode: 400,
				code:       "InvalidParameter",
				message:    "invalid",
			},
			permanent: true,
		},
		{
			name: "unauthorized 401 is permanent",
			err: fakeServiceError{
				statusCode: 401,
				code:       "NotAuthenticated",
				message:    "auth",
			},
			permanent: true,
		},
		{
			name: "forbidden 403 is permanent",
			err: fakeServiceError{
				statusCode: 403,
				code:       "NotAuthorized",
				message:    "forbidden",
			},
			permanent: true,
		},
		{
			name: "too many requests 429 is retryable",
			err: fakeServiceError{
				statusCode: 429,
				code:       "TooManyRequests",
				message:    "throttled",
			},
			permanent: false,
		},
		{
			name: "server error 500 is retryable",
			err: fakeServiceError{
				statusCode: 500,
				code:       "InternalServerError",
				message:    "server",
			},
			permanent: false,
		},
		{
			name:      "non service error is retryable",
			err:       errors.New("network timeout"),
			permanent: false,
		},
		{
			name:      "wrapped service error still detected",
			err:       fmt.Errorf("wrapper: %w", fakeServiceError{statusCode: 404, code: "NotFound", message: "missing"}),
			permanent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.permanent, isPermanentMonitoringError(tt.err))
		})
	}
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

func TestNewMonitoringClientFromHost(t *testing.T) {
	authID := component.MustNewID("oraclecloudauth")
	cfg := &Config{
		Region: "us-phoenix-1",
		Auth: configoptional.Some(configauth.Config{
			AuthenticatorID: authID,
		}),
	}

	client, err := newMonitoringClientFromHost(cfg, zap.NewNop(), fakeHost{
		extensions: map[component.ID]component.Component{
			authID: &fakeOracleCloudAuthExtension{provider: common.DefaultConfigProvider()},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, client)
}
