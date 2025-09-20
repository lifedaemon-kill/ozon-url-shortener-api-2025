package config

import (
	"errors"
	"main/internal/errs"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Yaml YamlConfig
	Env  EnvConfig
}

type YamlConfig struct {
}

type EnvConfig struct {
}

func Load(yamlPath, envPath string) (*Config, error) {
	yaml := YamlConfig{}
	if err := cleanenv.ReadConfig(yamlPath, yaml); err != nil {
		return nil, errors.Join(errs.ErrorYamlConfigFileNotFound, err)
	}
	env := EnvConfig{}
	if err := cleanenv.ReadConfig(envPath, env); err != nil {
		return nil, errors.Join(errs.ErrorEnvConfigFileNotFound, err)
	}
	return &Config{Yaml: yaml, Env: env}, nil
}
