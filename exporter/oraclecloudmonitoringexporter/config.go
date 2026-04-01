package oraclecloudmonitoringexporter

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	monitoringCompartmentIdKey    string = "oracle_cloud.monitoring.compartment.id" // Reserved keys for routing
	monitoringNamespaceKey        string = "oracle_cloud.monitoring.namespace" // Reserved keys for routing
	defaultMaxPastTimestampSkew          = 2 * time.Hour
	defaultMaxFutureTimestampSkew        = 10 * time.Minute
)

var reservedRoutingKeys = [...]string{monitoringCompartmentIdKey, monitoringNamespaceKey}

type Config struct {
	Region                 string                                     `mapstructure:"region"`
	Auth                   configoptional.Optional[configauth.Config] `mapstructure:"auth"`
	CompartmentId          string                                     `mapstructure:"compartment_id"`
	Namespace              string                                     `mapstructure:"namespace"`
	MaxPastTimestampSkew   time.Duration                              `mapstructure:"max_past_timestamp_skew"`
	MaxFutureTimestampSkew time.Duration                              `mapstructure:"max_future_timestamp_skew"`

	TimeoutSettings exporterhelper.TimeoutConfig                             `mapstructure:"timeout"`
	RetrySettings   configretry.BackOffConfig                                `mapstructure:"retry_on_failure"`
	QueueSettings   configoptional.Optional[exporterhelper.QueueBatchConfig] `mapstructure:"sending_queue"`
}

func (c *Config) Validate() error {
	if c.Region == "" {
		return errors.New(`"region" must be set`)
	}
	if !c.Auth.HasValue() {
		return errors.New(`"auth.authenticator" must be set and point to oraclecloudauthextension`)
	}
	if (c.CompartmentId == "") != (c.Namespace == "") {
		return errors.New(`"compartment_id" and "namespace" must be set together for exporter routing fallback`)
	}
	if c.MaxPastTimestampSkew < 0 {
		return errors.New(`"max_past_timestamp_skew" must be non-negative`)
	}
	if c.MaxPastTimestampSkew > defaultMaxPastTimestampSkew {
		return fmt.Errorf(`"max_past_timestamp_skew" must not exceed %s`, defaultMaxPastTimestampSkew)
	}
	if c.MaxFutureTimestampSkew < 0 {
		return errors.New(`"max_future_timestamp_skew" must be non-negative`)
	}
	if c.MaxFutureTimestampSkew > defaultMaxFutureTimestampSkew {
		return fmt.Errorf(`"max_future_timestamp_skew" must not exceed %s`, defaultMaxFutureTimestampSkew)
	}
	return nil
}

func (c *Config) effectiveMaxPastTimestampSkew() time.Duration {
	if c == nil || c.MaxPastTimestampSkew == 0 {
		return defaultMaxPastTimestampSkew
	}
	return c.MaxPastTimestampSkew
}

func (c *Config) effectiveMaxFutureTimestampSkew() time.Duration {
	if c == nil || c.MaxFutureTimestampSkew == 0 {
		return defaultMaxFutureTimestampSkew
	}
	return c.MaxFutureTimestampSkew
}

func isReservedAttributeKey(k string) bool {
	for _, rKey := range reservedRoutingKeys {
		if k == rKey {
			return true
		}
	}
	return false
}
