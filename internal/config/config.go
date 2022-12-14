/*
 * Copyright © 2022 Durudex
 *
 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Default config path.
const defaultConfigPath string = "configs/main"

type (
	// Config variables.
	Config struct {
		GRPC     GRPCConfig     `mapstructure:"grpc"`
		Database DatabaseConfig `mapstructure:"database"`
		Code     CodeConfig     `mapstructure:"code"`
		Service  ServiceConfig  `mapstructure:"service"`
	}

	// gRPC server config variables.
	GRPCConfig struct {
		Host string    `mapstructure:"host"`
		Port string    `mapstructure:"port"`
		TLS  TLSConfig `mapstructure:"tls"`
	}

	// TLS config variables.
	TLSConfig struct {
		Enable bool   `mapstructure:"enable"`
		CACert string `mapstructure:"ca-cert"`
		Cert   string `mapstructure:"cert"`
		Key    string `mapstructure:"key"`
	}

	// Database config variables.
	DatabaseConfig struct{ Redis RedisConfig }

	// Redis config variables.
	RedisConfig struct{ URL string }

	// Code config variables.
	CodeConfig struct {
		TTL       time.Duration `mapstructure:"ttl"`
		MaxLength int64         `mapstructure:"max-length"`
		MinLength int64         `mapstructure:"min-length"`
	}

	// Service base config.
	Service struct {
		Addr string    `mapstructure:"addr"`
		TLS  TLSConfig `mapstructure:"tls"`
	}

	// Services config variables.
	ServiceConfig struct {
		Email Service `mapstructure:"email"`
	}
)

// Creating a new config.
func NewConfig() (*Config, error) {
	log.Debug().Msg("Creating a new config...")

	// Parsing specified when starting the config file.
	if err := parseConfigFile(); err != nil {
		return nil, err
	}

	var cfg Config

	// Unmarshal config keys.
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Set configurations from environment.
	setFromEnv(&cfg)

	return &cfg, nil
}

// Parsing specified when starting the config file.
func parseConfigFile() error {
	// Get config path variable.
	configPath := os.Getenv("CONFIG_PATH")

	// Check is config path variable empty.
	if configPath == "" {
		configPath = defaultConfigPath
	}

	log.Debug().Msgf("Parsing config file: %s", configPath)

	// Split path to folder and file.
	dir, file := filepath.Split(configPath)

	viper.AddConfigPath(dir)
	viper.SetConfigName(file)

	// Read config file.
	return viper.ReadInConfig()
}

// Set configurations from environment.
func setFromEnv(cfg *Config) {
	log.Debug().Msg("Set configurations from environment...")

	// Redis database configurations.
	cfg.Database.Redis.URL = os.Getenv("REDIS_URL")
}
