package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"boilerplate/internal/app/boilerplate"
	"boilerplate/internal/model"
	"boilerplate/internal/pkg/version"
)

const (
	envPrefix = "BOILERPLATE"
)

func main() {
	versionInfo := version.Get()
	fmt.Println("Git version:", versionInfo.GitVersion)
	fmt.Println("Build date:", versionInfo.BuildDate)
	fmt.Println("Git commit:", versionInfo.GitCommit)
	fmt.Println("Platform:", versionInfo.Platform)

	var config model.Config
	rootCmd := &cobra.Command{
		Use:     "app",
		Short:   "App Short Description",
		Long:    "App Long Description",
		Example: "app",
		Run: func(_ *cobra.Command, _ []string) {
			app := boilerplate.New(&config)
			err := app.Run()
			if err != nil {
				panic(err)
			}
		},
	}

	// Log
	rootCmd.PersistentFlags().StringVarP(&config.LogLevel, "log-level", "", "debug", "Log Level")

	// DB
	rootCmd.PersistentFlags().StringVarP(&config.DB.Host, "db-host", "", "localhost", "Database Host")
	rootCmd.PersistentFlags().StringVarP(&config.DB.Port, "db-port", "", "5432", "Database Port")
	rootCmd.PersistentFlags().StringVarP(&config.DB.User, "db-user", "", "postgres", "Database User")
	rootCmd.PersistentFlags().StringVarP(&config.DB.Password, "db-password", "", "postgres", "Database Password")
	rootCmd.PersistentFlags().StringVarP(&config.DB.Name, "db-name", "", "boilerplate", "Database Name")
	rootCmd.PersistentFlags().BoolVarP(&config.DB.SslMode, "db-ssl-mode", "", false, "Database SSL Mode")

	// API
	rootCmd.PersistentFlags().StringVarP(&config.API.Port, "api-port", "", "8080", "API Port")
	rootCmd.PersistentFlags().StringVarP(&config.API.AccessPrivateKey, "api-access-private-key", "", "dd4dcf2eae3c3a6f097d69f49ce584852d66ac85505f5d264e1b6fb8f90d9019", "API Access Private Key")
	rootCmd.PersistentFlags().IntVarP(&config.API.AccessTokenTTL, "api-access-token-ttl", "", 600, "API Access Token TTL")
	rootCmd.PersistentFlags().IntVarP(&config.API.RefreshTokenTTL, "api-refresh-token-ttl", "", 2592000, "API Refresh Token TTL")

	cobra.OnInitialize(func() {
		v := viper.New()
		v.SetEnvPrefix(envPrefix)
	})

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
