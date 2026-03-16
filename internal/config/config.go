package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type DeployEntry struct {
	Path      string `yaml:"path"`
	Repo      string `yaml:"repo"`
	DeployKey string `yaml:"deploy_key"`
	Script    string `yaml:"script"`
}

type ConfigFile struct {
	DeployEntries []DeployEntry `yaml:"deploy"`
}

func LoadConfigFile(cfgPath string) (*ConfigFile, error) {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	cfg := &ConfigFile{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Expand ~ to home directory in paths
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil // continue without expansion
	}
	for i := range cfg.DeployEntries {
		if strings.HasPrefix(cfg.DeployEntries[i].Path, "~/") {
			cfg.DeployEntries[i].Path = filepath.Join(home, cfg.DeployEntries[i].Path[2:])
		}
		if strings.HasPrefix(cfg.DeployEntries[i].DeployKey, "~/") {
			cfg.DeployEntries[i].DeployKey = filepath.Join(home, cfg.DeployEntries[i].DeployKey[2:])
		}
	}

	return cfg, nil
}
