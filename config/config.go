package config

import "encoding/json"

type Config struct {
	DeployEnv string
	Listen    string
	Gateway struct {
		Endpoint string
	}
	Debug struct {
		Trace            bool
		GRPCTraceAddress string
	}
	Authorize struct {
		RawTokens []string
	}
	MetaInfoRedis struct {
		Addr     string
		Password string
		DB       int
	}
}

func (c *Config) String() string {
	b, e := json.Marshal(c)
	if e != nil {
		return "{}"
	}
	return string(b)
}

var (
	config = &Config{}
)

func GetConfig() *Config {
	return config
}
