package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type RotationConfig struct {
	AwsIamUsers []AwsIamUser `yaml:"aws_iam_users,omitempty"`
	Notifiers   []string     `yaml:"notifiers,omitempty"`
	SafeMode    bool         `yaml:"safe_mode,omitempty"`
}

type AwsIamUser struct {
	Username string                   `yaml:"username,omitempty"`
	Days     int                      `yaml:"days,omitempty"`
	Targets  []map[string]interface{} `yaml:"targets,omitempty"`
}

func LoadFromYamlFile(filePath string) (*RotationConfig, error) {
	config := RotationConfig{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.UnmarshalStrict(data, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
