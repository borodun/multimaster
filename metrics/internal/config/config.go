package config

type Database struct {
	URL  string `yaml:"url,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type Spec struct {
	Databases             []Database `yaml:"databases,omitempty"`
	ListenPort            int        `yaml:"listenPort,omitempty"`
	Interval              int        `yaml:"interval,omitempty"`
	QueryTimeout          int        `yaml:"queryTimeout,omitempty"`
	ConnectionPoolMaxSize int        `yaml:"connectionPoolMaxSize,omitempty"`
}

type Config struct {
	Spec Spec `yaml:"spec,omitempty"`
}
