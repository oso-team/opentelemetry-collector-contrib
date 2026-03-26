package oraclecloudmonitoringexporter

import (
	"fmt"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/monitoring"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func translateMetrics(md pmetric.Metrics, compIdCfg string, nsCfg string) ([]monitoring.MetricDataDetails, int, error) {
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
						detail, err := buildMetricDataDetails(m.Name(), resourceAttrs, dp.Attributes(), compIdCfg, nsCfg, numberDataPointValue(dp), dataPointTimestamp(dp.Timestamp()))
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
						detail, err := buildMetricDataDetails(m.Name(), resourceAttrs, dp.Attributes(), compIdCfg, nsCfg, numberDataPointValue(dp), dataPointTimestamp(dp.Timestamp()))
						if err != nil {
							return err
						}
						out = append(out, detail)
						return nil
					}); err != nil {
						return nil, dropped, err
					}
				case pmetric.MetricTypeHistogram:
					if err := forEachHistogramDataPoint(m.Histogram().DataPoints(), func(dp pmetric.HistogramDataPoint) error {
						// For histogram datapoints, emit count always and sum when present.
						countDetail, err := buildMetricDataDetails(m.Name()+".count", resourceAttrs, dp.Attributes(), compIdCfg, nsCfg, float64(dp.Count()), dataPointTimestamp(dp.Timestamp()))
						if err != nil {
							return err
						}
						out = append(out, countDetail)

						if dp.HasSum() {
							sumDetail, err := buildMetricDataDetails(m.Name()+".sum", resourceAttrs, dp.Attributes(), compIdCfg, nsCfg, dp.Sum(), dataPointTimestamp(dp.Timestamp()))
							if err != nil {
								return err
							}
							out = append(out, sumDetail)
						}
						return nil
					}); err != nil {
						return nil, dropped, err
					}
				case pmetric.MetricTypeExponentialHistogram:
					if err := forEachExponentialHistogramDataPoint(m.ExponentialHistogram().DataPoints(), func(dp pmetric.ExponentialHistogramDataPoint) error {
						// For exponential histogram datapoints, emit count always and sum when present.
						countDetail, err := buildMetricDataDetails(m.Name()+".count", resourceAttrs, dp.Attributes(), compIdCfg, nsCfg, float64(dp.Count()), dataPointTimestamp(dp.Timestamp()))
						if err != nil {
							return err
						}
						out = append(out, countDetail)

						if dp.HasSum() {
							sumDetail, err := buildMetricDataDetails(m.Name()+".sum", resourceAttrs, dp.Attributes(), compIdCfg, nsCfg, dp.Sum(), dataPointTimestamp(dp.Timestamp()))
							if err != nil {
								return err
							}
							out = append(out, sumDetail)
						}
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

func forEachHistogramDataPoint(points pmetric.HistogramDataPointSlice, fn func(dp pmetric.HistogramDataPoint) error) error {
	for i := 0; i < points.Len(); i++ {
		if err := fn(points.At(i)); err != nil {
			return err
		}
	}
	return nil
}

func forEachExponentialHistogramDataPoint(points pmetric.ExponentialHistogramDataPointSlice, fn func(dp pmetric.ExponentialHistogramDataPoint) error) error {
	for i := 0; i < points.Len(); i++ {
		if err := fn(points.At(i)); err != nil {
			return err
		}
	}
	return nil
}

func buildMetricDataDetails(metricName string, resourceAttrs pcommon.Map, datapointAttrs pcommon.Map, compIdCfg string, nsCfg string, value float64, ts time.Time) (monitoring.MetricDataDetails, error) {
	compartmentID, namespace, err := getRoutingValues(resourceAttrs, compIdCfg, nsCfg)
	if err != nil {
		return monitoring.MetricDataDetails{}, err
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

func getRoutingValues(resAttr pcommon.Map, compIdCfg string, nsCfg string) (string, string, error) {
	compartmentId := ""
	namespace := ""
	if c, ok := resAttr.Get(monitoringCompartmentIdKey); ok {
		compartmentId = c.AsString()
	}
	if n, ok := resAttr.Get(monitoringNamespaceKey); ok {
		namespace = n.AsString()
	}

	// Routing key precedence:
	// 1) use resource attrs only when both values are present and non-empty
	// 2) otherwise use exporter config fallback only when both values are present and non-empty
	if compartmentId != "" && namespace != "" {
		return compartmentId, namespace, nil
	}
	if compIdCfg != "" && nsCfg != "" {
		return compIdCfg, nsCfg, nil
	}

	return "", "", fmt.Errorf("missing required attributes %q and %q", monitoringCompartmentIdKey, monitoringNamespaceKey)
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
