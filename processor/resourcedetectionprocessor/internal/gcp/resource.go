// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package gcp // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/gcp"

import (
	"fmt"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp"

	localMetadata "github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor/internal/gcp/internal/metadata"
)

// These helpers live here rather than in the generated internal/metadata
// package: they depend on the GCP detector library, and a hand-written file
// in the metadata package would make the generated config schema inherit
// that dependency. They only use exported ResourceBuilder methods.

func setFromCallable(set func(string), detect func() (string, error)) error {
	v, err := detect()
	if err != nil {
		return err
	}
	set(v)
	return nil
}

func setZoneAndRegion(rb *localMetadata.ResourceBuilder, detect func() (string, string, error)) error {
	zone, region, err := detect()
	if err != nil {
		return err
	}
	rb.SetCloudAvailabilityZone(zone)
	rb.SetCloudRegion(region)
	return nil
}

func setZoneOrRegion(rb *localMetadata.ResourceBuilder, detect func() (string, gcp.LocationType, error)) error {
	v, locType, err := detect()
	if err != nil {
		return err
	}
	switch locType {
	case gcp.Zone:
		rb.SetCloudAvailabilityZone(v)
	case gcp.Region:
		rb.SetCloudRegion(v)
	default:
		return fmt.Errorf("location must be zone or region. Got %v", locType)
	}
	return nil
}

func setManagedInstanceGroup(rb *localMetadata.ResourceBuilder, detect func() (gcp.ManagedInstanceGroup, error)) error {
	v, err := detect()
	if err != nil {
		return err
	}
	if v.Name != "" {
		rb.SetGcpGceInstanceGroupManagerName(v.Name)
	}
	switch v.Type {
	case gcp.Zone:
		rb.SetGcpGceInstanceGroupManagerZone(v.Location)
	case gcp.Region:
		rb.SetGcpGceInstanceGroupManagerRegion(v.Location)
	}
	return nil
}
