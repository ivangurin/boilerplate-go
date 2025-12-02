package suite_provider

import "boilerplate/internal/model"

func GetConfig() *model.Config {
	return &model.Config{
		DbDsn: "postgres://postgres:postgres@localhost:5432/boilerplate_test",
	}
}
