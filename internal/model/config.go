package model

import "fmt"

type Config struct {
	LogLevel string    `yaml:"log_level" json:"log_level"`
	DB       ConfigDB  `yaml:"db" json:"db"`
	API      ConfigAPI `yaml:"api" json:"api"`
}

type ConfigDB struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Name     string `yaml:"name" json:"name"`
	SslMode  bool   `yaml:"ssl-mode" json:"ssl-mode"`
}

type ConfigAPI struct {
	Port             string `yaml:"port" json:"port"`
	AccessPrivateKey string `yaml:"access-private-key" json:"access-private-key"`
	AccessTokenTTL   int    `yaml:"access-token-ttl" json:"token-ttl"`
	RefreshTokenTTL  int    `yaml:"refresh-token-ttl" json:"refresh-token-ttl"`
}

func (c ConfigDB) GetDSN() string {
	sslMode := "disable"
	if c.SslMode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, sslMode)

	return dsn
}
