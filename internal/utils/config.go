package utils

import "time"

type Config struct {
	Address        string        `env:"ADDRESS,required" envDefault:"127.0.0.1:8080" json:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s" json:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s" json:"REPORT_INTERVAL"`
}

func NewConfig() *Config {
	return &Config{}
}
