package model

import "fmt"

type Config struct {
	LogLevel string     `yaml:"log_level" json:"log_level" mapstructure:"log-level" validate:"required"`
	DB       ConfigDB   `yaml:"db" json:"db" mapstructure:"db"`
	API      ConfigAPI  `yaml:"api" json:"api" mapstructure:"api"`
	S3       ConfigS3   `yaml:"s3" json:"s3" mapstructure:"s3" validate:"required"`
	Nats     ConfigNats `yaml:"nats" json:"nats" mapstructure:"nats" validate:"required"`
	Mail     ConfigMail `yaml:"mail" json:"mail" mapstructure:"mail" validate:"required"`
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
	HTTPPort         string `yaml:"http-port" json:"http-port" mapstructure:"http-port" validate:"required"`
	GRPCPort         string `yaml:"grpc-port" json:"grpc-port" mapstructure:"grpc-port" validate:"required"`
	AccessPrivateKey string `yaml:"access-private-key" json:"access-private-key" mapstructure:"access-private-key" validate:"required"`
	AccessTokenTTL   int    `yaml:"access-token-ttl" json:"token-ttl" mapstructure:"access-token-ttl" validate:"required"`
	RefreshTokenTTL  int    `yaml:"refresh-token-ttl" json:"refresh-token-ttl" mapstructure:"refresh-token-ttl" validate:"required"`
}

type ConfigS3 struct {
	Host      string `yaml:"host" json:"host" mapstructure:"host" validate:"required"`
	Port      string `yaml:"port" json:"port" mapstructure:"port" validate:"required"`
	AccessKey string `yaml:"access-key" json:"access-key" mapstructure:"access-key" validate:"required"`
	SecretKey string `yaml:"secret-key" json:"secret-key" mapstructure:"secret-key" validate:"required"`
	Bucket    string `yaml:"bucket" json:"bucket" mapstructure:"bucket" validate:"required"`
}

type ConfigNats struct {
	Host     string `yaml:"host" json:"host" mapstructure:"host" validate:"required"`
	Port     string `yaml:"port" json:"port" mapstructure:"port" validate:"required"`
	HTTPPort string `yaml:"http-port" json:"http-port" mapstructure:"http-port" validate:"required"`
	Domain   string `yaml:"domain" json:"domain" mapstructure:"domain" validate:"required"`
	DataDir  string `yaml:"data-dir" json:"data-dir" mapstructure:"data-dir" validate:"required"`
}

type ConfigMail struct {
	SMTPHost string `yaml:"smtp-host" json:"smtp-host" mapstructure:"smtp-host" validate:"required"`
	SMTPPort string `yaml:"smtp-port" json:"smtp-port" mapstructure:"smtp-port" validate:"required"`
	Username string `yaml:"username" json:"username" mapstructure:"username" validate:"required"`
	Password string `yaml:"password" json:"password" mapstructure:"password" validate:"required"`
	SSL      bool   `yaml:"ssl" json:"ssl" mapstructure:"ssl"`
	TLS      bool   `yaml:"tls" json:"tls" mapstructure:"tls"`
	From     string `yaml:"from" json:"from" mapstructure:"from" validate:"required,email"`
}

func (c ConfigDB) GetDSN() string {
	sslMode := "disable"
	if c.SslMode {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, sslMode)

	return dsn
}
