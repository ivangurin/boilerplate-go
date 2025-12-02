package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"boilerplate/internal/app/boilerplate"
	"boilerplate/internal/model"
)

func main() {
	var config model.Config
	rootCmd := &cobra.Command{
		Use:     "app",
		Short:   "App Short Description",
		Long:    "App Long Description",
		Example: "app",
		Run: func(_ *cobra.Command, _ []string) {
			app := boilerplate.New(config)
			err := app.Run()
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&config.DbDsn, "db-dsn", "", "postgres://postgres:postgres@localhost:5432/boilerplate", "Database DSN")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
