package oraclecloudmonitoringexporter

import (
	"errors"

	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	Region string                                     `mapstructure:"region"`
	Auth   configoptional.Optional[configauth.Config] `mapstructure:"auth"`

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
	return nil
}

func isReservedAttributeKey(k string) bool {
	switch k {
	case "monitoring_compartment_id", "oci.monitoring.compartment.id", "monitoring_namespace", "oci.monitoring.namespace":
		return true
	default:
		return false
	}
}
