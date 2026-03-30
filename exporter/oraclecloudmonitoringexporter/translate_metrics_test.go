package oraclecloudmonitoringexporter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

const testTimestamp = pcommon.Timestamp(1710000000000000000)

func TestTranslateMetricsGauge(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	rm.Resource().Attributes().PutStr("service.name", "checkout")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetDoubleValue(12.5)
	dp.Attributes().PutStr("host.name", "node-1")

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)

	require.Equal(t, "otel_demo", *data[0].Namespace)
	require.Equal(t, "ocid1.compartment.oc1..aaaa", *data[0].CompartmentId)
	require.Equal(t, "cpu.utilization", *data[0].Name)
	require.Equal(t, "checkout", data[0].Dimensions["service.name"])
	require.Equal(t, "node-1", data[0].Dimensions["host.name"])
	_, hasNamespace := data[0].Dimensions["oracle_cloud.monitoring.namespace"]
	_, hasCompartment := data[0].Dimensions["oracle_cloud.monitoring.compartment.id"]
	require.False(t, hasNamespace)
	require.False(t, hasCompartment)
}

func TestTranslateMetricsMissingAttributes(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetDoubleValue(1)

	_, _, err := translateMetrics(md, "", "")
	require.Error(t, err)
}

func TestTranslateMetricsHistogramCountAndSum(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	rm.Resource().Attributes().PutStr("service.name", "checkout")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("latency")
	h := m.SetEmptyHistogram()
	dp := h.DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetCount(5)
	dp.SetSum(42.5)
	dp.Attributes().PutStr("host.name", "node-1")

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 2)

	require.Equal(t, "latency.count", *data[0].Name)
	require.Equal(t, 5.0, *data[0].Datapoints[0].Value)
	require.Equal(t, "checkout", data[0].Dimensions["service.name"])
	require.Equal(t, "node-1", data[0].Dimensions["host.name"])

	require.Equal(t, "latency.sum", *data[1].Name)
	require.Equal(t, 42.5, *data[1].Datapoints[0].Value)
	require.Equal(t, "checkout", data[1].Dimensions["service.name"])
	require.Equal(t, "node-1", data[1].Dimensions["host.name"])
}

func TestTranslateMetricsHistogramCountOnlyWhenSumMissing(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("request.duration")
	h := m.SetEmptyHistogram()
	dp := h.DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetCount(7)
	dp.Attributes().PutStr("host.name", "node-1")

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)
	require.Equal(t, "request.duration.count", *data[0].Name)
	require.Equal(t, 7.0, *data[0].Datapoints[0].Value)
}

func TestTranslateMetricsExponentialHistogramCountAndSum(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	rm.Resource().Attributes().PutStr("service.name", "checkout")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("exp.latency")
	eh := m.SetEmptyExponentialHistogram()
	dp := eh.DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetCount(9)
	dp.SetSum(100.5)
	dp.Attributes().PutStr("host.name", "node-2")

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 2)

	require.Equal(t, "exp.latency.count", *data[0].Name)
	require.Equal(t, 9.0, *data[0].Datapoints[0].Value)
	require.Equal(t, "checkout", data[0].Dimensions["service.name"])
	require.Equal(t, "node-2", data[0].Dimensions["host.name"])

	require.Equal(t, "exp.latency.sum", *data[1].Name)
	require.Equal(t, 100.5, *data[1].Datapoints[0].Value)
	require.Equal(t, "checkout", data[1].Dimensions["service.name"])
	require.Equal(t, "node-2", data[1].Dimensions["host.name"])
}

func TestTranslateMetricsExponentialHistogramCountOnlyWhenSumMissing(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("exp.request.duration")
	eh := m.SetEmptyExponentialHistogram()
	dp := eh.DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetCount(4)
	dp.Attributes().PutStr("host.name", "node-2")

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)
	require.Equal(t, "exp.request.duration.count", *data[0].Name)
	require.Equal(t, 4.0, *data[0].Datapoints[0].Value)
}

func TestTranslateMetricsDropsUnsupported(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("summary_latency")
	s := m.SetEmptySummary()
	s.DataPoints().AppendEmpty()

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Empty(t, data)
	require.Equal(t, 1, dropped)
}

func TestTranslateMetricsLegacyMonitoringKeysNotReserved(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	rm.Resource().Attributes().PutStr("monitoring_namespace", "legacy-ns")
	rm.Resource().Attributes().PutStr("monitoring_compartment_id", "legacy-compartment")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetDoubleValue(12.5)

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)
	require.Equal(t, "legacy-ns", data[0].Dimensions["monitoring_namespace"])
	require.Equal(t, "legacy-compartment", data[0].Dimensions["monitoring_compartment_id"])
}

func TestTranslateMetricsDropsOverLimitDimensions(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetDoubleValue(12.5)
	for i := 0; i < maxDimensionsPerMetric+1; i++ {
		dp.Attributes().PutStr(fmt.Sprintf("dim.%02d", i), "value")
	}

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Empty(t, data)
	require.Equal(t, 1, dropped)
}

func TestTranslateMetricsFallsBackToConfigRouting(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("service.name", "checkout")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetDoubleValue(12.5)

	data, dropped, err := translateMetrics(md, "ocid1.compartment.oc1..cfg", "cfg_ns")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)
	require.Equal(t, "cfg_ns", *data[0].Namespace)
	require.Equal(t, "ocid1.compartment.oc1..cfg", *data[0].CompartmentId)
}

func TestTranslateMetricsResourceRoutingOverridesConfig(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..resource")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "resource_ns")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetDoubleValue(1)
	dp.Attributes().PutStr("host.name", "node-1")

	data, dropped, err := translateMetrics(md, "ocid1.compartment.oc1..cfg", "cfg_ns")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)
	require.Equal(t, "resource_ns", *data[0].Namespace)
	require.Equal(t, "ocid1.compartment.oc1..resource", *data[0].CompartmentId)
}

func TestTranslateMetricsUsesConfigWhenResourceRoutingIsPartial(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..resource")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(testTimestamp)
	dp.SetDoubleValue(1)
	dp.Attributes().PutStr("host.name", "node-1")

	data, dropped, err := translateMetrics(md, "ocid1.compartment.oc1..cfg", "cfg_ns")
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)
	require.Equal(t, "cfg_ns", *data[0].Namespace)
	require.Equal(t, "ocid1.compartment.oc1..cfg", *data[0].CompartmentId)
}

func TestTranslateMetricsDropsMissingTimestamp(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oracle_cloud.monitoring.namespace", "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetDoubleValue(10)

	data, dropped, err := translateMetrics(md, "", "")
	require.NoError(t, err)
	require.Empty(t, data)
	require.Equal(t, 1, dropped)
}
