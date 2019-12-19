package memviz

type Config struct{}

type Configurator func(*Config)

func defaultConfig() *Config {
	return &Config{}
}

func New(configurators ...Configurator) *Config {
	config := defaultConfig()
	for _, configurator := range configurators {
		configurator(config)
	}
	return config
}
