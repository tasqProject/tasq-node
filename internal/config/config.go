// Package config loads node settings from the environment, with flags in main
// taking precedence.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/tasqProject/tasq-node/internal/offer"
)

// Config is the resolved node configuration.
type Config struct {
	Endpoint string       // coordinator base URL
	GPU      string       // GPU model to advertise
	Region   string       // region to advertise
	Price    string       // price per hour in credits
	Modes    []offer.Mode // modes to offer
	KeyPath  string       // path to the operator key
}

// FromEnv reads configuration from TASQ_* environment variables, falling back to
// sensible defaults. The serve command lets flags override these.
func FromEnv() Config {
	return Config{
		Endpoint: env("TASQ_ENDPOINT", "https://api.tasqnetwork.io"),
		GPU:      env("TASQ_GPU", ""),
		Region:   env("TASQ_REGION", ""),
		Price:    env("TASQ_PRICE", ""),
		Modes:    parseModes(env("TASQ_MODES", "R")),
		KeyPath:  env("TASQ_KEY", "operator.key"),
	}
}

// Validate checks the configuration is complete enough to start serving.
func (c Config) Validate() error {
	if c.GPU == "" || c.Region == "" || c.Price == "" {
		return fmt.Errorf("config: gpu, region and price are required")
	}
	if len(c.Modes) == 0 {
		return fmt.Errorf("config: at least one mode is required")
	}
	return nil
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func parseModes(s string) []offer.Mode {
	out := make([]offer.Mode, 0, len(s))
	for _, c := range strings.ToUpper(s) {
		switch c {
		case 'A':
			out = append(out, offer.Attested)
		case 'R':
			out = append(out, offer.Redundant)
		case 'P':
			out = append(out, offer.Proven)
		}
	}
	return out
}
