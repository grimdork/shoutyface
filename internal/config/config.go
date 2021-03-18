package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config for the main server.
type Config struct {
	// Address of server in the format ip:port.
	Address string `json:"address"`
	// Logs holds log paths, if outputting to files.
	Logs Logs `json:"logs"`
}

// Logs define files to save output to.
type Logs struct {
	Messages string `json:"messages"`
	Errors   string `json:"errors"`
}

// Default returns a reasonable collection of defaults.
func Default() *Config {
	cfg := Config{
		Address: "127.0.0.1:15015",
		Logs: Logs{
			Messages: "messages.log",
			Errors:   "errors.log",
		},
	}
	return &cfg
}

// LoadConfig attempts to load the JSON file at the provided path.
func Load(path string) (*Config, error) {
	cfg := Default()
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, nil
	}

	err = json.Unmarshal(in, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Save to file.
func (cfg *Config) Save(path string) error {
	var b []byte
	var err error
	b, err = json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	return err
}
