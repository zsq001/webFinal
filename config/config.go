/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type SqliteDBConfig struct {
	Path string `yaml:"path"`
}

type Config struct {
	BindAddr string         `yaml:"bind_addr"`
	Database SqliteDBConfig `yaml:"database"`
	JwtKey   string         `yaml:"jwt_key"`
}

var CONFIG = Config{
	BindAddr: "localhost:10086",
	Database: SqliteDBConfig{
		Path: "./db.sqlite3",
	},
	JwtKey: "1145141919810",
}

func ReadConfig(filename string) error {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &CONFIG)
	if err != nil {
		return err
	}
	return nil
}
