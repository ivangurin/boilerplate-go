package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
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
	log.Println("Build date:", versionInfo.BuildDate)
	log.Println("Platform:", versionInfo.Platform)
	log.Println("Git version:", versionInfo.GitVersion)
	log.Println("Git commit:", versionInfo.GitCommit)
	fmt.Println()

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AddConfigPath(".")

	config := &model.Config{}
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

	if err := defineFlags(rootCmd, config); err != nil {
		log.Printf("Error defining flags: %v", err)
		os.Exit(1)
	}

	rootCmd.PreRun = func(_ *cobra.Command, _ []string) {
		configFile, err := rootCmd.PersistentFlags().GetString("config")
		if err != nil {
			log.Printf("Error getting config flag: %v", err)
			os.Exit(1)
		}

		if configFile != "" {
			viper.SetConfigFile(configFile)
			if err := viper.ReadInConfig(); err != nil {
				log.Printf("Error reading config file: %v\n", err)
				os.Exit(1)
			}

			log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
		}

		if err := viper.Unmarshal(config); err != nil {
			log.Printf("Error unmarshaling config: %v\n", err)
			os.Exit(1)
		}

		if err := validateConfig(config); err != nil {
			log.Printf("Config validation error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error executing root command: %v", err)
		os.Exit(1)
	}
}

func defineFlags(cmd *cobra.Command, config *model.Config) error {
	var err error

	// Config
	cmd.PersistentFlags().String("config", "", "Path to config file(Ex: ./config.yaml)")

	// Log
	if err = bindStringVar(cmd, &config.LogLevel, "log-level", "debug", "Log Level"); err != nil {
		return fmt.Errorf("bind log-level: %w", err)
	}

	// DB
	if err = bindStringVar(cmd, &config.DB.Host, "db.host", "localhost", "Database Host"); err != nil {
		return fmt.Errorf("bind db.host: %w", err)
	}
	if err = bindStringVar(cmd, &config.DB.Port, "db.port", "5432", "Database Port"); err != nil {
		return fmt.Errorf("bind db.port: %w", err)
	}
	if err = bindStringVar(cmd, &config.DB.User, "db.user", "postgres", "Database User"); err != nil {
		return fmt.Errorf("bind db.user: %w", err)
	}
	if err = bindStringVar(cmd, &config.DB.Password, "db.password", "postgres", "Database Password"); err != nil {
		return fmt.Errorf("bind db.password: %w", err)
	}
	if err = bindStringVar(cmd, &config.DB.Name, "db.name", "boilerplate", "Database Name"); err != nil {
		return fmt.Errorf("bind db.name: %w", err)
	}
	if err = bindBoolVar(cmd, &config.DB.SslMode, "db.ssl-mode", false, "Database SSL Mode"); err != nil {
		return fmt.Errorf("bind db.ssl-mode: %w", err)
	}

	// API
	if err = bindStringVar(cmd, &config.API.Host, "api.host", "127.0.0.1", "API Host"); err != nil {
		return fmt.Errorf("bind api.host: %w", err)
	}
	if err = bindStringVar(cmd, &config.API.HTTPPort, "api.http-port", "8080", "API HTTP Port"); err != nil {
		return fmt.Errorf("bind api.http-port: %w", err)
	}
	if err = bindStringVar(cmd, &config.API.GRPCPort, "api.grpc-port", "8082", "API GRPC Port"); err != nil {
		return fmt.Errorf("bind api.grpc-port: %w", err)
	}
	if err = bindStringVar(cmd, &config.API.AccessPrivateKey, "api.access-private-key", "dd4dcf2eae3c3a6f097d69f49ce584852d66ac85505f5d264e1b6fb8f90d9019", "API Access Private Key"); err != nil {
		return fmt.Errorf("bind api.access-private-key: %w", err)
	}
	if err = bindIntVar(cmd, &config.API.AccessTokenTTL, "api.access-token-ttl", 600, "API Access Token TTL"); err != nil {
		return fmt.Errorf("bind api.access-token-ttl: %w", err)
	}
	if err = bindIntVar(cmd, &config.API.RefreshTokenTTL, "api.refresh-token-ttl", 2592000, "API Refresh Token TTL"); err != nil {
		return fmt.Errorf("bind api.refresh-token-ttl: %w", err)
	}

	return nil
}

func bindStringVar(cmd *cobra.Command, p *string, name string, value string, usage string) error {
	cmd.PersistentFlags().StringVar(p, name, value, usage)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		return fmt.Errorf("bind flag %s: %w", name, err)
	}
	return nil
}

func bindIntVar(cmd *cobra.Command, p *int, name string, value int, usage string) error {
	cmd.PersistentFlags().IntVar(p, name, value, usage)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		return fmt.Errorf("bind flag %s: %w", name, err)
	}
	return nil
}

func bindBoolVar(cmd *cobra.Command, p *bool, name string, value bool, usage string) error {
	cmd.PersistentFlags().BoolVar(p, name, value, usage)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		return fmt.Errorf("bind flag %s: %w", name, err)
	}
	return nil
}

func validateConfig(config *model.Config) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				return fmt.Errorf("field '%s' validation failed: %s", e.Field(), e.Tag())
			}
		}
		return err
	}
	return nil
}
