package service

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

// CatalogConfig represents details re: which modules' services should be
// included or excluded from the catalog
type CatalogConfig struct {
	MinStability Stability
	UseV2Guid    bool
}

type tempCatalogConfig struct {
	CatalogConfig
	MinStabilityStr string `envconfig:"MIN_STABILITY" default:"STABLE"`
	UseV2GuidStr    string `envconfig:"USE_V2_GUID" default:"false"`
}

// NewCatalogConfigWithDefaults returns a CatalogConfig object with default
// values already applied. Callers are then free to set custom values for the
// remaining fields and/or override default values.
func NewCatalogConfigWithDefaults() CatalogConfig {
	return CatalogConfig{
		MinStability: StabilityPreview,
		UseV2Guid:    false,
	}
}

// GetCatalogConfigFromEnvironment returns catalog configuration
func GetCatalogConfigFromEnvironment() (CatalogConfig, error) {
	c := tempCatalogConfig{
		CatalogConfig: NewCatalogConfigWithDefaults(),
	}
	err := envconfig.Process("", &c)
	if err != nil {
		return c.CatalogConfig, err
	}
	minStabilityStr := strings.ToUpper(c.MinStabilityStr)
	switch minStabilityStr {
	case "EXPERIMENTAL":
		c.MinStability = StabilityExperimental
	case "PREVIEW":
		c.MinStability = StabilityPreview
	case "STABLE":
		c.MinStability = StabilityStable
	default:
		return c.CatalogConfig, fmt.Errorf(
			`unrecognized stability level "%s"`,
			minStabilityStr,
		)
	}
	useV2GuidStr := strings.ToUpper(c.UseV2GuidStr)
	switch useV2GuidStr {
	case "TRUE":
		c.UseV2Guid = true
	case "FALSE":
		c.UseV2Guid = false
	default:
		return c.CatalogConfig, fmt.Errorf(
			`unrecognized useV2Guid boolean "%s"`,
			useV2GuidStr,
		)
	}
	return c.CatalogConfig, nil
}
