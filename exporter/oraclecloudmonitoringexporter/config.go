package oraclecloudmonitoringexporter

import (
	"errors"

	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Reserved keys for routing
const (
	monitoringCompartmentIdKey string = "oracle_cloud.monitoring.compartment.id"
	monitoringNamespaceKey     string = "oracle_cloud.monitoring.namespace"
)

var reservedRoutingKeys = [...]string{monitoringCompartmentIdKey, monitoringNamespaceKey}

type Config struct {
	Region        string                                     `mapstructure:"region"`
	Auth          configoptional.Optional[configauth.Config] `mapstructure:"auth"`
	CompartmentId string                                     `mapstructure:"compartment_id"`
	Namespace     string                                     `mapstructure:"namespace"`

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
	return nil
}

func isReservedAttributeKey(k string) bool {
	for _, rKey := range reservedRoutingKeys {
		if k == rKey {
			return true
		}
	}
	return false
}
