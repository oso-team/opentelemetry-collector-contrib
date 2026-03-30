package oraclecloudmonitoringexporter

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/oracle/oci-go-sdk/v65/monitoring"
)

const (
	maxDimensionsPerMetric     = 20
	maxUniqueStreamsPerRequest = 50
)

type batchedMetrics struct {
	batches [][]monitoring.MetricDataDetails
}

func buildMetricBatches(metricData []monitoring.MetricDataDetails) batchedMetrics {
	if len(metricData) == 0 {
		return batchedMetrics{}
	}

	// streamGroups is a map of stream key to corresponding metricDataDetails list
	streamGroups := make(map[string][]monitoring.MetricDataDetails, len(metricData))
	// collection of stream keys i.e, all unique streams.
	streamKeys := make([]string, 0, len(metricData))

	for _, metric := range metricData {
		key := buildStreamKey(metric)
		if _, found := streamGroups[key]; !found {
			streamKeys = append(streamKeys, key)
		}
		streamGroups[key] = append(streamGroups[key], metric)
	}

	if len(streamKeys) == 0 {
		return batchedMetrics{}
	}

	batches := make([][]monitoring.MetricDataDetails, 0, intCeil(len(streamKeys), maxUniqueStreamsPerRequest))
	currentBatch := make([]monitoring.MetricDataDetails, 0)
	currUniqueStreams := 0

	for _, key := range streamKeys {
		group := streamGroups[key]
		if currUniqueStreams == maxUniqueStreamsPerRequest {
			batches = append(batches, currentBatch)
			currentBatch = make([]monitoring.MetricDataDetails, 0)
			currUniqueStreams = 0
		}

		currentBatch = append(currentBatch, group...)
		currUniqueStreams++
	}

	if len(currentBatch) > 0 {
		batches = append(batches, currentBatch)
	}

	return batchedMetrics{batches: batches}
}

func buildStreamKey(metric monitoring.MetricDataDetails) string {
	return fmt.Sprintf("%s|%s|%s|%s",
		*metric.CompartmentId,
		*metric.Namespace,
		*metric.Name,
		normalizeDimensions(metric.Dimensions),
	)
}

func normalizeDimensions(dimensions map[string]string) string {
	if len(dimensions) == 0 {
		return ""
	}

	// slices.Sorted extracts and sorts keys in one step
	keys := slices.Sorted(maps.Keys(dimensions))

	var builder strings.Builder
	for i, key := range keys {
		if i > 0 {
			builder.WriteString(",")
		}
		fmt.Fprintf(&builder, "%s=%s", key, dimensions[key])
	}
	return builder.String()
}

func isDimensionsWithinLimit(metric monitoring.MetricDataDetails) bool {
	count := len(metric.Dimensions)
	return count > 0 && count <= maxDimensionsPerMetric
}

func intCeil(num1, num2 int) int {
	return (num1 + num2 - 1) / num2
}
