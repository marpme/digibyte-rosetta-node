// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package configuration

import (
	"github.com/pkg/errors"
	uconfig "go.uber.org/config"
)

const (
	// ConfigPath is the ENV variable that will be looked for if not existent
	ConfigPath = "ROSETTA_DGB_CONFIG_PATH"
)

type (
	// NetworkIdentifier specifies a blockchain network it will run on
	NetworkIdentifier struct {
		Blockchain string `yaml:"blockchain"`
		Network    string `yaml:"network"`
	}
	// Currency is the specification for the given currency
	Currency struct {
		Symbol   string `yaml:"symbol"`
		Decimals int32  `yaml:"decimals"`
	}
	// Server represents setting for this rosetta data server
	Server struct {
		Port string `yaml:"port"`
	}

	// Node specifies the connection details towards a given node
	Node struct {
		Endpoint   string `yaml:"endpoint"`
		TLSEnabled bool   `yaml:"tlsEnabled"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
	}

	// Redis specifies the connection details towards a given node
	Redis struct {
		Addr     string `yaml:"address"`
		Password string `yaml:"password"`
	}

	// Version is representing the version of the specification
	// for both rosetta and the given node
	Version struct {
		RosettaVersion  string `yaml:"rosettaVersion"`
		DigibyteVersion string `yaml:"digibyteVersion"`
	}

	// Config is the overall configuration for this service layer
	Config struct {
		NetworkIdentifier NetworkIdentifier `yaml:"network_identifier"`
		Currency          Currency          `yaml:"currency"`
		Server            Server            `yaml:"server"`
		Node              Node              `yaml:"node"`
		Version           Version           `yaml:"version"`
		Redis             Redis             `yaml:"redis"`
	}
)

// New parses a new config file and creates everything necessary
func New(path string) (cfg *Config, err error) {
	opts := []uconfig.YAMLOption{uconfig.File(path)}
	yaml, err := uconfig.NewYAML(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}
	cfg = &Config{}
	if err := yaml.Get(uconfig.Root).Populate(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal YAML config to struct")
	}
	return
}
