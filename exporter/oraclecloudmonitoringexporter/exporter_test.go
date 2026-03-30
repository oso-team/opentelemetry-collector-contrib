package oraclecloudmonitoringexporter

import (
	"errors"
	"fmt"
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
				err: errors.New("transport failure"),
			},
		},
	}

	md := metricWithRoutingAttrs()
	err := exp.pushMetricsData(t.Context(), md)
	require.Error(t, err)
	require.True(t, consumererror.IsPermanent(err))
}

func TestPushMetricsDataSendErrorIsPermanent(t *testing.T) {
	exp := &metricsExporter{
		cfg:    &Config{},
		logger: zap.NewNop(),
		client: &oracleCloudMonitoringClient{
			logger: zap.NewNop(),
			client: &fakeMonitoringClient{
				err: errors.New("sdk retries exhausted"),
			},
		},
	}

	md := metricWithRoutingAttrs()
	err := exp.pushMetricsData(t.Context(), md)
	require.Error(t, err)
	require.True(t, consumererror.IsPermanent(err))
}

func TestPushMetricsDataSplitsPlannedRequestsSequentially(t *testing.T) {
	fake := &fakeMonitoringClient{}
	exp := &metricsExporter{
		cfg:    &Config{},
		logger: zap.NewNop(),
		client: &oracleCloudMonitoringClient{
			logger: zap.NewNop(),
			client: fake,
		},
	}

	md := metricWithManyUniqueStreams(maxUniqueStreamsPerRequest + 1)
	err := exp.pushMetricsData(t.Context(), md)
	require.NoError(t, err)
	require.Len(t, fake.requests, 2)
	require.Len(t, fake.requests[0].PostMetricDataDetails.MetricData, maxUniqueStreamsPerRequest)
	require.Len(t, fake.requests[1].PostMetricDataDetails.MetricData, 1)
	require.NotNil(t, fake.requests[0].RequestMetadata.RetryPolicy)
	require.NotNil(t, fake.requests[1].RequestMetadata.RetryPolicy)
}

func TestPushMetricsDataStopsAfterFirstFailedPlannedRequest(t *testing.T) {
	fake := &fakeMonitoringClient{
		errByCall: map[int]error{
			2: errors.New("second request failed"),
		},
	}
	exp := &metricsExporter{
		cfg:    &Config{},
		logger: zap.NewNop(),
		client: &oracleCloudMonitoringClient{
			logger: zap.NewNop(),
			client: fake,
		},
	}

	md := metricWithManyUniqueStreams(maxUniqueStreamsPerRequest*2 + 1)
	err := exp.pushMetricsData(t.Context(), md)
	require.Error(t, err)
	require.True(t, consumererror.IsPermanent(err))
	require.Len(t, fake.requests, 2)
}

func TestPushMetricsDataCountsDimensionLimitedMetricsAsTranslationDrops(t *testing.T) {
	fake := &fakeMonitoringClient{}
	exp := &metricsExporter{
		cfg:    &Config{},
		logger: zap.NewNop(),
		client: &oracleCloudMonitoringClient{
			logger: zap.NewNop(),
			client: fake,
		},
	}

	md := metricWithDimensionCount(maxDimensionsPerMetric + 1)
	err := exp.pushMetricsData(t.Context(), md)
	require.NoError(t, err)
	require.Empty(t, fake.requests)
}

func metricWithRoutingAttrs() pmetric.Metrics {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr(monitoringCompartmentIdKey, "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr(monitoringNamespaceKey, "otel_demo")
	rm.Resource().Attributes().PutStr("service.name", "checkout")
	sm := rm.ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(pcommon.Timestamp(1710000000000000000))
	dp.SetDoubleValue(1)
	return md
}

func metricWithManyUniqueStreams(streams int) pmetric.Metrics {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr(monitoringCompartmentIdKey, "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr(monitoringNamespaceKey, "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()

	for i := 0; i < streams; i++ {
		m := sm.Metrics().AppendEmpty()
		m.SetName("cpu.utilization")
		dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
		dp.SetTimestamp(pcommon.Timestamp(1710000000000000000))
		dp.SetDoubleValue(1)
		dp.Attributes().PutStr("host.name", fmt.Sprintf("node-%03d", i))
		dp.Attributes().PutStr("stream.id", fmt.Sprintf("%03d", i))
	}

	return md
}

func metricWithDimensionCount(count int) pmetric.Metrics {
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
	for i := 0; i < count; i++ {
		dp.Attributes().PutStr(fmt.Sprintf("dim.%02d", i), "value")
	}
	return md
}
