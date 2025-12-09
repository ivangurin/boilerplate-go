package model

import "fmt"

type Config struct {
	LogLevel string    `yaml:"log_level" json:"log_level" mapstructure:"log-level" validate:"required"`
	DB       ConfigDB  `yaml:"db" json:"db" mapstructure:"db"`
	API      ConfigAPI `yaml:"api" json:"api" mapstructure:"api"`
}

type ConfigDB struct {
	Host     string `yaml:"host" json:"host" mapstructure:"host" validate:"required"`
	Port     string `yaml:"port" json:"port" mapstructure:"port" validate:"required"`
	User     string `yaml:"user" json:"user" mapstructure:"user" validate:"required"`
	Password string `yaml:"password" json:"password" mapstructure:"password" validate:"required"`
	Name     string `yaml:"name" json:"name" mapstructure:"name" validate:"required"`
	SslMode  bool   `yaml:"ssl-mode" json:"ssl-mode" mapstructure:"ssl-mode"`
}

type ConfigAPI struct {
	Host             string `yaml:"host" json:"host" mapstructure:"host" validate:"required"`
	Port             string `yaml:"port" json:"port" mapstructure:"port" validate:"required"`
	AccessPrivateKey string `yaml:"access-private-key" json:"access-private-key" mapstructure:"access-private-key" validate:"required"`
	AccessTokenTTL   int    `yaml:"access-token-ttl" json:"token-ttl" mapstructure:"access-token-ttl" validate:"required"`
	RefreshTokenTTL  int    `yaml:"refresh-token-ttl" json:"refresh-token-ttl" mapstructure:"refresh-token-ttl" validate:"required"`
}

func (c ConfigDB) GetDSN() string {
	sslMode := "disable"
	if c.SslMode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, sslMode)

	return dsn
}
