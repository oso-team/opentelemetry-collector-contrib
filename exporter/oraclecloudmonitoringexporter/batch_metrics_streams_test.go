package oraclecloudmonitoringexporter

import (
	"fmt"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"github.com/stretchr/testify/require"
)

func TestPlanMetricBatchesDropsZeroAndOverLimitDimensions(t *testing.T) {
	valid := testMetricDataDetail("metric.valid", map[string]string{"service.name": "checkout"})

	planned := buildMetricBatches([]monitoring.MetricDataDetails{valid})

	require.Len(t, planned.batches, 1)
	require.Len(t, planned.batches[0], 1)
	require.Equal(t, "metric.valid", *planned.batches[0][0].Name)
}

func TestPlanMetricBatchesCountsReservedRoutingKeysOutsideDimensions(t *testing.T) {
	metric := testMetricDataDetail("metric.valid", map[string]string{"service.name": "checkout"})

	planned := buildMetricBatches([]monitoring.MetricDataDetails{metric})

	require.Len(t, planned.batches, 1)
	require.Len(t, planned.batches[0], 1)
	require.Equal(t, 1, len(planned.batches[0][0].Dimensions))
}

func TestPlanMetricBatchesKeepsSameStreamTogether(t *testing.T) {
	first := testMetricDataDetail("metric.requests", map[string]string{"service.name": "checkout"})
	second := testMetricDataDetail("metric.requests", map[string]string{"service.name": "checkout"})
	other := testMetricDataDetail("metric.other", map[string]string{"service.name": "checkout"})

	planned := buildMetricBatches([]monitoring.MetricDataDetails{first, second, other})

	require.Len(t, planned.batches, 1)
	require.Len(t, planned.batches[0], 3)
	require.Equal(t, *first.Name, *planned.batches[0][0].Name)
	require.Equal(t, *second.Name, *planned.batches[0][1].Name)
	require.Equal(t, *other.Name, *planned.batches[0][2].Name)
}

func TestPlanMetricBatchesSplitsAfterFiftyUniqueStreams(t *testing.T) {
	metrics := make([]monitoring.MetricDataDetails, 0, maxUniqueStreamsPerRequest+1)
	for i := 0; i < maxUniqueStreamsPerRequest+1; i++ {
		metrics = append(metrics, testMetricDataDetail(
			fmt.Sprintf("metric.%02d", i),
			map[string]string{"service.name": fmt.Sprintf("checkout-%02d", i)},
		))
	}

	planned := buildMetricBatches(metrics)

	require.Len(t, planned.batches, 2)
	require.Len(t, planned.batches[0], maxUniqueStreamsPerRequest)
	require.Len(t, planned.batches[1], 1)
}

func TestPlanMetricBatchesSupportsMixedRoutingValues(t *testing.T) {
	first := testMetricDataDetail("metric.cpu", map[string]string{"service.name": "checkout"})
	second := testMetricDataDetail("metric.cpu", map[string]string{"service.name": "checkout"})
	namespace := "other_ns"
	second.Namespace = common.String(namespace)

	planned := buildMetricBatches([]monitoring.MetricDataDetails{first, second})

	require.Len(t, planned.batches, 1)
	require.Len(t, planned.batches[0], 2)
	require.NotEqual(t, buildStreamKey(planned.batches[0][0]), buildStreamKey(planned.batches[0][1]))
}

func testMetricDataDetail(name string, dimensions map[string]string) monitoring.MetricDataDetails {
	namespace := "otel_demo"
	compartmentID := "ocid1.compartment.oc1..aaaa"
	return monitoring.MetricDataDetails{
		Name:          common.String(name),
		Namespace:     common.String(namespace),
		CompartmentId: common.String(compartmentID),
		Dimensions:    dimensions,
	}
}
