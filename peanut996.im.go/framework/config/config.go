package config

type Config struct {
	Mongo `yaml:"mongo"`
	Redis `yaml:"redis"`
}

type Mongo struct {
	Url string `yaml:"url"`
}

type Redis struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Passwd string `yaml:"passwd"`
	DB     int    `yaml:"db"`
}
