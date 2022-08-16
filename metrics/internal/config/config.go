package config

type Database struct {
	URL    string            `yaml:"url,omitempty"`
	Name   string            `yaml:"name,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
}

type Metrics struct {
	Databases             []Database        `yaml:"databases,omitempty"`
	ListenPort            int               `yaml:"listenPort,omitempty"`
	Interval              int               `yaml:"interval,omitempty"`
	QueryTimeout          int               `yaml:"queryTimeout,omitempty"`
	ConnectionPoolMaxSize int               `yaml:"connectionPoolMaxSize,omitempty"`
	Labels                map[string]string `yaml:"labels,omitempty"`
}

type Config struct {
	Metrics Metrics `yaml:"metrics,omitempty"`
}
