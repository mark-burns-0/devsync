package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type SyncConfig struct {
	ProjectsRoot string `env:"PROJECTS_ROOT" env-default:""`
	GitDir       string `env:"GIT_DIR" env-default:".git"`
}

func New(path string) (*SyncConfig, error) {
	cfg := &SyncConfig{}

	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
