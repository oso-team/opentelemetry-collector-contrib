package oraclecloudmonitoringexporter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestTranslateMetricsGauge(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oci.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oci.monitoring.namespace", "otel_demo")
	rm.Resource().Attributes().PutStr("service.name", "checkout")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetTimestamp(1710000000000000000)
	dp.SetDoubleValue(12.5)
	dp.Attributes().PutStr("host.name", "node-1")

	data, dropped, err := translateMetrics(md)
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)

	require.Equal(t, "otel_demo", *data[0].Namespace)
	require.Equal(t, "ocid1.compartment.oc1..aaaa", *data[0].CompartmentId)
	require.Equal(t, "cpu.utilization", *data[0].Name)
	require.Equal(t, "checkout", data[0].Dimensions["service.name"])
	require.Equal(t, "node-1", data[0].Dimensions["host.name"])
	_, hasNamespace := data[0].Dimensions["oci.monitoring.namespace"]
	_, hasCompartment := data[0].Dimensions["oci.monitoring.compartment.id"]
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
	dp.SetDoubleValue(1)

	_, _, err := translateMetrics(md)
	require.Error(t, err)
}

func TestTranslateMetricsDropsUnsupported(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oci.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oci.monitoring.namespace", "otel_demo")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("latency")
	h := m.SetEmptyHistogram()
	h.DataPoints().AppendEmpty()

	data, dropped, err := translateMetrics(md)
	require.NoError(t, err)
	require.Empty(t, data)
	require.Equal(t, 1, dropped)
}

func TestTranslateMetricsLegacyMonitoringKeysNotReserved(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("oci.monitoring.compartment.id", "ocid1.compartment.oc1..aaaa")
	rm.Resource().Attributes().PutStr("oci.monitoring.namespace", "otel_demo")
	rm.Resource().Attributes().PutStr("monitoring_namespace", "legacy-ns")
	rm.Resource().Attributes().PutStr("monitoring_compartment_id", "legacy-compartment")
	sm := rm.ScopeMetrics().AppendEmpty()

	m := sm.Metrics().AppendEmpty()
	m.SetName("cpu.utilization")
	dp := m.SetEmptyGauge().DataPoints().AppendEmpty()
	dp.SetDoubleValue(12.5)

	data, dropped, err := translateMetrics(md)
	require.NoError(t, err)
	require.Equal(t, 0, dropped)
	require.Len(t, data, 1)
	require.Equal(t, "legacy-ns", data[0].Dimensions["monitoring_namespace"])
	require.Equal(t, "legacy-compartment", data[0].Dimensions["monitoring_compartment_id"])
}
