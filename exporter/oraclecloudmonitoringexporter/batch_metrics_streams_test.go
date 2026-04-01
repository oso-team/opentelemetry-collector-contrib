package oraclecloudmonitoringexporter

import (
	"fmt"
	"strings"
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

func TestIsMetricDimensionsValid(t *testing.T) {
	tests := []struct {
		name       string
		dimensions map[string]string
		want       bool
	}{
		{
			name:       "valid_dimensions",
			dimensions: map[string]string{"service.name": "checkout"},
			want:       true,
		},
		{
			name:       "empty_key",
			dimensions: map[string]string{"": "checkout"},
			want:       false,
		},
		{
			name:       "key_with_space",
			dimensions: map[string]string{"service name": "checkout"},
			want:       false,
		},
		{
			name:       "key_with_non_ascii",
			dimensions: map[string]string{"service.namé": "checkout"},
			want:       false,
		},
		{
			name:       "key_too_long",
			dimensions: map[string]string{strings.Repeat("a", maxDimensionKeyLength+1): "checkout"},
			want:       false,
		},
		{
			name:       "key_at_limit",
			dimensions: map[string]string{strings.Repeat("a", maxDimensionKeyLength): "checkout"},
			want:       true,
		},
		{
			name:       "empty_value",
			dimensions: map[string]string{"service.name": ""},
			want:       false,
		},
		{
			name:       "value_too_long",
			dimensions: map[string]string{"service.name": strings.Repeat("a", maxDimensionValueLength+1)},
			want:       false,
		},
		{
			name:       "value_at_limit",
			dimensions: map[string]string{"service.name": strings.Repeat("a", maxDimensionValueLength)},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isMetricDimensionsValid(testMetricDataDetail("metric.cpu", tt.dimensions)))
		})
	}
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
