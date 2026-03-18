package oraclecloudmonitoringexporter

import (
	"errors"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var errMissingRequiredAttributes = errors.New(`missing required attributes "monitoring_compartment_id" (or "oci.monitoring.compartment.id") and "monitoring_namespace" (or "oci.monitoring.namespace")`)

func translateMetrics(md pmetric.Metrics) ([]monitoring.MetricDataDetails, int, error) {
	out := make([]monitoring.MetricDataDetails, 0)
	dropped := 0

	resourceMetrics := md.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		rm := resourceMetrics.At(i)
		resourceAttrs := rm.Resource().Attributes()

		scopeMetrics := rm.ScopeMetrics()
		for j := 0; j < scopeMetrics.Len(); j++ {
			metrics := scopeMetrics.At(j).Metrics()
			for k := 0; k < metrics.Len(); k++ {
				m := metrics.At(k)
				switch m.Type() {
				case pmetric.MetricTypeGauge:
					if err := forEachNumberDataPoint(m.Gauge().DataPoints(), func(dp pmetric.NumberDataPoint) error {
						detail, err := buildMetricDataDetails(m.Name(), resourceAttrs, dp.Attributes(), numberDataPointValue(dp), dataPointTimestamp(dp.Timestamp()))
						if err != nil {
							return err
						}
						out = append(out, detail)
						return nil
					}); err != nil {
						return nil, dropped, err
					}
				case pmetric.MetricTypeSum:
					if err := forEachNumberDataPoint(m.Sum().DataPoints(), func(dp pmetric.NumberDataPoint) error {
						detail, err := buildMetricDataDetails(m.Name(), resourceAttrs, dp.Attributes(), numberDataPointValue(dp), dataPointTimestamp(dp.Timestamp()))
						if err != nil {
							return err
						}
						out = append(out, detail)
						return nil
					}); err != nil {
						return nil, dropped, err
					}
				default:
					dropped += estimateDroppedCount(m)
				}
			}
		}
	}
	return out, dropped, nil
}

func forEachNumberDataPoint(points pmetric.NumberDataPointSlice, fn func(dp pmetric.NumberDataPoint) error) error {
	for i := 0; i < points.Len(); i++ {
		if err := fn(points.At(i)); err != nil {
			return err
		}
	}
	return nil
}

func buildMetricDataDetails(metricName string, resourceAttrs pcommon.Map, datapointAttrs pcommon.Map, value float64, ts time.Time) (monitoring.MetricDataDetails, error) {
	compartmentID := getAttributeString(datapointAttrs, resourceAttrs, "monitoring_compartment_id", "oci.monitoring.compartment.id")
	namespace := getAttributeString(datapointAttrs, resourceAttrs, "monitoring_namespace", "oci.monitoring.namespace")
	if compartmentID == "" || namespace == "" {
		return monitoring.MetricDataDetails{}, errMissingRequiredAttributes
	}

	dimensions := mergedDimensions(resourceAttrs, datapointAttrs)

	return monitoring.MetricDataDetails{
		Namespace:     common.String(namespace),
		CompartmentId: common.String(compartmentID),
		Name:          common.String(metricName),
		Dimensions:    dimensions,
		Datapoints: []monitoring.Datapoint{{
			Timestamp: &common.SDKTime{Time: ts},
			Value:     common.Float64(value),
		}},
	}, nil
}

func mergedDimensions(resourceAttrs pcommon.Map, datapointAttrs pcommon.Map) map[string]string {
	dimensions := make(map[string]string, resourceAttrs.Len()+datapointAttrs.Len())
	resourceAttrs.Range(func(k string, v pcommon.Value) bool {
		if !isReservedAttributeKey(k) {
			dimensions[k] = v.AsString()
		}
		return true
	})
	datapointAttrs.Range(func(k string, v pcommon.Value) bool {
		if !isReservedAttributeKey(k) {
			dimensions[k] = v.AsString()
		}
		return true
	})
	return dimensions
}

func getAttributeString(primary pcommon.Map, fallback pcommon.Map, keys ...string) string {
	for _, key := range keys {
		if value, ok := primary.Get(key); ok {
			return value.AsString()
		}
		if value, ok := fallback.Get(key); ok {
			return value.AsString()
		}
	}
	return ""
}

func numberDataPointValue(dp pmetric.NumberDataPoint) float64 {
	switch dp.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return float64(dp.IntValue())
	default:
		return dp.DoubleValue()
	}
}

func dataPointTimestamp(ts pcommon.Timestamp) time.Time {
	if ts == 0 {
		return time.Now().UTC()
	}
	return time.Unix(0, int64(ts)).UTC()
}

func estimateDroppedCount(metric pmetric.Metric) int {
	switch metric.Type() {
	case pmetric.MetricTypeHistogram:
		return metric.Histogram().DataPoints().Len()
	case pmetric.MetricTypeExponentialHistogram:
		return metric.ExponentialHistogram().DataPoints().Len()
	case pmetric.MetricTypeSummary:
		return metric.Summary().DataPoints().Len()
	default:
		return 0
	}
}
