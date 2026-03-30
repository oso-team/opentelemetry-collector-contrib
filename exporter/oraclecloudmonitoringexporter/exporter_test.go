package oraclecloudmonitoringexporter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

func TestPushMetricsDataTranslateErrorIsPermanent(t *testing.T) {
	exp := &metricsExporter{
		cfg:    &Config{},
		logger: zap.NewNop(),
		client: &oracleCloudMonitoringClient{
			logger: zap.NewNop(),
			client: &fakeMonitoringClient{},
		},
	}

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(pcommon.Timestamp(1710000000000000000))
	dp.SetDoubleValue(1)

	err := exp.pushMetricsData(t.Context(), md)
	require.Error(t, err)
	require.True(t, consumererror.IsPermanent(err))
}

func TestPushMetricsDataSend400IsPermanent(t *testing.T) {
	exp := &metricsExporter{
		cfg:    &Config{},
		logger: zap.NewNop(),
		client: &oracleCloudMonitoringClient{
			logger: zap.NewNop(),
			client: &fakeMonitoringClient{
				err: fakeServiceError{statusCode: 400, code: "InvalidParameter", message: "invalid"},
			},
		},
	}

	md := metricWithRoutingAttrs()
	err := exp.pushMetricsData(t.Context(), md)
	require.Error(t, err)
	require.True(t, consumererror.IsPermanent(err))
}

func TestPushMetricsDataSend429IsRetryable(t *testing.T) {
	exp := &metricsExporter{
		cfg:    &Config{},
		logger: zap.NewNop(),
		client: &oracleCloudMonitoringClient{
			logger: zap.NewNop(),
			client: &fakeMonitoringClient{
				err: fakeServiceError{statusCode: 429, code: "TooManyRequests", message: "throttled"},
			},
		},
	}

	md := metricWithRoutingAttrs()
	err := exp.pushMetricsData(t.Context(), md)
	require.Error(t, err)
	require.False(t, consumererror.IsPermanent(err))
}

func metricWithRoutingAttrs() pmetric.Metrics {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr(monitoringCompartmentIdKey, "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr(monitoringNamespaceKey, "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(pcommon.Timestamp(1710000000000000000))
	dp.SetDoubleValue(1)
	return md
}
