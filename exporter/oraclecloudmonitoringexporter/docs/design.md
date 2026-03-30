# Design (Draft)

## OTel to Oracle Cloud Monitoring - Metric type mapping

***Oracle Cloud Monitoring metrics schema:***

**compartmentId** - Mandatory to scope

**namespace** - Mandatory to scope

**name** - Metric name (Mandatory)

**dimensions** - Analogous to labels (Mandatory)

**datapoints** - Sample with timestamp, value and count. (Mandatory)

**resourceGroup** - Custom string that you can match when retrieving custom metrics. Only one resource group can be applied per metric. (_Optional_)

**metadata** - Properties describing metrics. These are not part of the unique fields identifying the metric. Example: `"unit": "bytes"` (_Optional_)

Example:
```json
[
  {
    "compartmentId": "ocid1.compartment.oc1..aaaaaaaan26oieeaintlbl6gptqxjdy74p7wp2lezetabcphzsf3f4xz32hq",
    "namespace": "sharhasstest",
    "name": "sample_metric",
    "dimensions": {
      "region": "us-phoenix-1",
      "app_name": "cli"
    },
    "datapoints": [
      {
        "timestamp": "2026-03-09T05:54:47Z",
        "value": 45.0,
        "count": 1
      }
    ]
  }
]
```

### Mappings

#### 1. **Counter** type mapping:
- Direct 1:1 mapping is possible
- Value from datapoint is mapped to value
```text
resource_metrics {
  resource {
    attributes { key: "service.name" value { string_value: "checkout" } }
    attributes { key: "oracle_cloud.monitoring.compartment.id" value { string_value: "ocid1.compartment.oc1..aaaa..." } }
    attributes { key: "oracle_cloud.monitoring.namespace" value { string_value: "otel_demo" } }
  }
  scope_metrics {
    metrics {
      name: "http.server.requests"
      unit: "{request}"
      sum {
        aggregation_temporality: AGGREGATION_TEMPORALITY_CUMULATIVE
        is_monotonic: true
        data_points {
          time_unix_nano: 1710000000000000000
          as_int: 120
          attributes { key: "http.method" value { string_value: "GET" } }
        }
      }
    }
  }
}
```
```json
{
  "metricData": [
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.requests",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 120
        }
      ]
    }
  ]
}
```
#### 2. **UpDownCounter** type mapping:
- Direct 1:1 mapping is possible
- Value from datapoint is maped to value
```text
resource_metrics {
  resource {
    attributes { key: "service.name" value { string_value: "checkout" } }
    attributes { key: "oracle_cloud.monitoring.compartment.id" value { string_value: "ocid1.compartment.oc1..aaaa..." } }
    attributes { key: "oracle_cloud.monitoring.namespace" value { string_value: "otel_demo" } }
  }
  scope_metrics {
    metrics {
      name: "active.sessions"
      unit: "{session}"
      sum {
        aggregation_temporality: AGGREGATION_TEMPORALITY_CUMULATIVE
        is_monotonic: false
        data_points {
          time_unix_nano: 1710000000000000000
          as_int: 42
          attributes { key: "env" value { string_value: "prod" } }
        }
      }
    }
  }
}
```
```json
{
  "metricData": [
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "active.sessions",
      "dimensions": {
        "service.name": "checkout",
        "env": "prod"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 42
        }
      ]
    }
  ]
}
```
#### 3. **Histogram** type mapping:
- We can either just use the scalar values of sum, count, min, max
- Or we can decompose the buckets
```text
resource_metrics {
  resource {
    attributes { key: "service.name" value { string_value: "checkout" } }
    attributes { key: "oracle_cloud.monitoring.compartment.id" value { string_value: "ocid1.compartment.oc1..aaaa..." } }
    attributes { key: "oracle_cloud.monitoring.namespace" value { string_value: "otel_demo" } }
  }
  scope_metrics {
    metrics {
      name: "http.server.duration"
      unit: "ms"
      histogram {
        aggregation_temporality: AGGREGATION_TEMPORALITY_CUMULATIVE
        data_points {
          time_unix_nano: 1710000000000000000
          count: 5
          sum: 220
          min: 10
          max: 95
          explicit_bounds: 10
          explicit_bounds: 25
          explicit_bounds: 50
          explicit_bounds: 100
          bucket_counts: 1
          bucket_counts: 2
          bucket_counts: 1
          bucket_counts: 1
          bucket_counts: 0
          attributes { key: "http.method" value { string_value: "GET" } }
        }
      }
    }
  }
}
```
```json
{
  "metricData": [
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 220
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_count",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 5
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_min",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 10
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_max",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 95
        }
      ]
    }
  ]
}
```
OR based on bucket boundary with cumulative sum.
```json
{
  "metricData": [
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_bucket",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET",
        "le": "10"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 1
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_bucket",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET",
        "le": "25"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 3
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_bucket",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET",
        "le": "50"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 4
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_bucket",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET",
        "le": "100"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 5
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_bucket",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET",
        "le": "+Inf"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 5
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_sum",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 220
        }
      ]
    },
    {
      "namespace": "otel_demo",
      "compartmentId": "ocid1.compartment.oc1..aaaa...",
      "name": "http.server.duration_count",
      "dimensions": {
        "service.name": "checkout",
        "http.method": "GET"
      },
      "datapoints": [
        {
          "timestamp": "2024-03-09T16:00:00Z",
          "value": 5
        }
      ]
    }
  ]
}
```
#### 4. **ExponentialHistograms** 
Can follow the same approach as Histogram
#### 5. **Summary** - Skip

## Data Validation and drop policy

- Missing routing attrs (`oracle_cloud.monitoring.compartment.id`/`oracle_cloud.monitoring.namespace`) 
    - Drop datapoint, count/log reason
- Unsupported metric type in current mode (Example: Summary)
    - Drop datapoint, count/log reason
- Malformed histogram (bucket count mismatch, invalid bounds)
    - Drop datapoint, count/log reason
- Missing timestamp
    - Drop datapoint, count/log reason
- Invalid attribute key/value for backend constraints
    ```
    - PostMetricData timestamp window: datapoints must be within 2 hours past and 10 minutes future of current time.
    - Per-call limits: max 20 dimensions per metric group, max 50 unique metric streams, 50 TPS per tenancy for post operation. (A metric group is the combination of a given metric, metric namespace, and tenancy for the purpose of determining limits.)
    - Dimension validation: key must be printable ASCII (no spaces), key max 256 chars, value max 512 chars, keys/values cannot be empty.
    ```
    Reference: [doc](https://docs.oracle.com/en-us/iaas/api/#/en/monitoring/20180401/MetricData/PostMetricData)
    - Drop offending datapoint or sanitize according to configured policy

## Rate Limit and Retry-Storm Handling

The exporter can rely on exporterhelper standard mechanisms:

- timeout
- retry with backoff
- sending queue

Expected behavior:

- Retry transient/network/5xx errors with backoff.
- Treat permanent validation/authz errors as non-retryable.
- Queue should be bounded to prevent memory issues.

