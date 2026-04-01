package oraclecloudmonitoringexporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configoptional"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid_with_auth_extension",
			cfg: Config{
				Region: "us-phoenix-1",
				Auth: configoptional.Some(configauth.Config{
					AuthenticatorID: component.MustNewID("oraclecloudauth"),
				}),
			},
		},
		{
			name:    "missing_region",
			cfg:     Config{},
			wantErr: true,
		},
		{
			name: "missing_auth",
			cfg: Config{
				Region: "us-phoenix-1",
			},
			wantErr: true,
		},
		{
			name: "fallback_config_both_values_set",
			cfg: Config{
				Region:        "us-phoenix-1",
				CompartmentId: "ocid1.compartment.oc1..aaaa",
				Namespace:     "otel_demo",
				Auth: configoptional.Some(configauth.Config{
					AuthenticatorID: component.MustNewID("oraclecloudauth"),
				}),
			},
		},
		{
			name: "fallback_config_missing_namespace",
			cfg: Config{
				Region:        "us-phoenix-1",
				CompartmentId: "ocid1.compartment.oc1..aaaa",
				Auth: configoptional.Some(configauth.Config{
					AuthenticatorID: component.MustNewID("oraclecloudauth"),
				}),
			},
			wantErr: true,
		},
		{
			name: "fallback_config_missing_compartment_id",
			cfg: Config{
				Region:    "us-phoenix-1",
				Namespace: "otel_demo",
				Auth: configoptional.Some(configauth.Config{
					AuthenticatorID: component.MustNewID("oraclecloudauth"),
				}),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestConfigValidateRejectsNegativeTimestampSkews(t *testing.T) {
	cfg := Config{
		Region: "us-phoenix-1",
		Auth: configoptional.Some(configauth.Config{
			AuthenticatorID: component.MustNewID("oraclecloudauth"),
		}),
		MaxPastTimestampSkew: -time.Second,
	}

	err := cfg.Validate()
	require.Error(t, err)
	require.ErrorContains(t, err, "max_past_timestamp_skew")

	cfg.MaxPastTimestampSkew = 0
	cfg.MaxFutureTimestampSkew = -time.Second

	err = cfg.Validate()
	require.Error(t, err)
	require.ErrorContains(t, err, "max_future_timestamp_skew")
}

func TestConfigValidateRejectsTimestampSkewsAboveOCILimits(t *testing.T) {
	cfg := Config{
		Region: "us-phoenix-1",
		Auth: configoptional.Some(configauth.Config{
			AuthenticatorID: component.MustNewID("oraclecloudauth"),
		}),
		MaxPastTimestampSkew: defaultMaxPastTimestampSkew + time.Second,
	}

	err := cfg.Validate()
	require.Error(t, err)
	require.ErrorContains(t, err, "max_past_timestamp_skew")

	cfg.MaxPastTimestampSkew = 0
	cfg.MaxFutureTimestampSkew = defaultMaxFutureTimestampSkew + time.Second

	err = cfg.Validate()
	require.Error(t, err)
	require.ErrorContains(t, err, "max_future_timestamp_skew")
}

func TestCreateDefaultConfigSetsTimestampSkewDefaults(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	require.Equal(t, defaultMaxPastTimestampSkew, cfg.MaxPastTimestampSkew)
	require.Equal(t, defaultMaxFutureTimestampSkew, cfg.MaxFutureTimestampSkew)
}

func TestConfigEffectiveTimestampSkewsUseDefaults(t *testing.T) {
	cfg := &Config{}

	require.Equal(t, defaultMaxPastTimestampSkew, cfg.effectiveMaxPastTimestampSkew())
	require.Equal(t, defaultMaxFutureTimestampSkew, cfg.effectiveMaxFutureTimestampSkew())

	cfg.MaxPastTimestampSkew = 30 * time.Minute
	cfg.MaxFutureTimestampSkew = 5 * time.Minute

	require.Equal(t, 30*time.Minute, cfg.effectiveMaxPastTimestampSkew())
	require.Equal(t, 5*time.Minute, cfg.effectiveMaxFutureTimestampSkew())
}
