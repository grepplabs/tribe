package config

import "github.com/rs/cors"

type CorsConfig struct {
	Enabled bool
	Options cors.Options
}
