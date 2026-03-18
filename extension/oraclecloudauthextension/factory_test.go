// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package oraclecloudauthextension

import (
	"testing"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/extension/extensiontest"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()
	require.NotNil(t, f)
	assert.Equal(t, "oraclecloudauth", f.Type().String())
}

func TestCreateExtension(t *testing.T) {
	origAPI := newAPIKeyProvider
	t.Cleanup(func() {
		newAPIKeyProvider = origAPI
	})
	newAPIKeyProvider = func(APIKeyConfig) (common.ConfigurationProvider, error) {
		return common.DefaultConfigProvider(), nil
	}

	cfg := createDefaultConfig()
	ext, err := createExtension(t.Context(), extensiontest.NewNopSettings(NewFactory().Type()), cfg)
	require.NoError(t, err)
	require.NotNil(t, ext)
}
