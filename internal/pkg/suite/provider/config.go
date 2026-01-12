package suite_provider

import "boilerplate/internal/model"

func (sp *Provider) GetConfig() *model.Config {
	if sp.config == nil {
		sp.config = InitConfig()
	}
	return sp.config
}

func InitConfig() *model.Config {
	return &model.Config{
		DB: model.ConfigDB{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Name:     "boilerplate_test",
			SslMode:  false,
		},
		API: model.ConfigAPI{
			HTTPPort:         "8080",
			GRPCPort:         "8082",
			AccessPrivateKey: "dd4dcf2eae3c3a6f097d69f49ce584852d66ac85505f5d264e1b6fb8f90d9019",
			AccessTokenTTL:   10,
			RefreshTokenTTL:  60,
		},
		S3: model.ConfigS3{
			Host:      "localhost",
			Port:      "9000",
			AccessKey: "admin",
			SecretKey: "password",
			Bucket:    "boilerplate-test",
		},
		Chrome: model.ConfigChrome{
			Host:    "localhost",
			Port:    "3000",
			Timeout: 30,
		},
	}
}
