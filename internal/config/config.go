package config

import "gopkg.in/yaml.v2"

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
	cfg := &ConfigFile{}

	yaml.Unmarshal([]byte(cfgPath), cfg)

	return cfg, nil
}
