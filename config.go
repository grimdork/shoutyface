package main

import (
	"os"
	"path/filepath"

	"github.com/francoispqt/gojay"
	"github.com/grimdork/xos"
)

// Config holds client configuration.
type Config struct {
	ServerURL  string
	AdminToken string
}

// UnmarshalJSONObject decodes this config from JSON via gojay.
func (cfg *Config) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "serverurl":
		return dec.String(&cfg.ServerURL)
	case "admintoken":
		return dec.String(&cfg.AdminToken)
	}
	return nil
}

// NKeys is required to unmarshal.
func (cfg *Config) NKeys() int {
	return 4
}

func (cfg *Config) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("admintoken", cfg.AdminToken)
	enc.StringKey("serverurl", cfg.ServerURL)
}

func (cfg *Config) IsNil() bool {
	return cfg == nil
}

func loadConfig() (*Config, error) {
	path, err := xos.NewConfig("shoutyface")
	if err != nil {
		return nil, err
	}

	var cfg Config
	f, err := os.Open(filepath.Join(path.Path(), "config.json"))
	if err != nil {
		return &cfg, nil
	}

	defer f.Close()
	dec := gojay.NewDecoder(f)
	err = dec.DecodeObject(&cfg)
	return &cfg, err
}

func saveConfig(cfg *Config) error {
	path, err := xos.NewConfig("shoutyface")
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Path(), 0700)
	if err != nil {
		return err
	}

	buf, err := gojay.MarshalJSONObject(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path.Path(), "config.json"), buf, 0600)
}
