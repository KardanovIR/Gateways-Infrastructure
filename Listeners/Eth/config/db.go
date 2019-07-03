package config

type DB struct {
	Name string `mapstructure:"NAME"`
	Host string `mapstructure:"HOST"`
}
