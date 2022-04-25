package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	"github.com/mohammadiahmad/QuicMessageBroker/internal/server"
)

type Config struct {
	Server *server.Config `mapstructure:"server"`
}

const (
	delimeter = "."
	seperator = "__"

	tagName = "mapstructure"

	upTemplate     = "================ Loaded Configuration ================"
	bottomTemplate = "======================================================"
)

func Load() (*Config, error) {
	k := koanf.New(delimeter)

	if err := LoadEnv(k); err != nil {
		return nil, fmt.Errorf("error loading default values: %v", err)
	}

	config := Config{}
	var tag = koanf.UnmarshalConf{Tag: tagName}
	if err := k.UnmarshalWithConf("", &config, tag); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %v", err)
	}

	// pretty print loaded configuration using provided template
	log.Printf("%s\n%v\n%s\n", upTemplate, spew.Sdump(config), bottomTemplate)
	return &config, nil
}

func LoadEnv(k *koanf.Koanf) error {
	var prefix = strings.ToUpper("quic") + seperator

	callback := func(source string) string {
		base := strings.ToLower(strings.TrimPrefix(source, prefix))
		return strings.ReplaceAll(base, seperator, delimeter)
	}

	// load environment variables
	if err := k.Load(env.Provider(prefix, delimeter, callback), nil); err != nil {
		return fmt.Errorf("error loading environment variables: %s", err)
	}

	return nil
}
