// Package conf handles the configuration of the applications. Yaml
// files are mapped with the struct
package conf

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config is the configuration struct. The config file config.yaml
// will unmarshaled to this struct.
type Config struct {
	DebugEnabled     bool
	Oauth2Enabled    bool
	ProfilingEnabled bool
	LogFlushInterval time.Duration
	URL              string
	RealURL          *url.URL //RealURL to our service endpoint parsed from URL
	AuthURL          string
	TokenURL         string
	Username         string //user to authenticate with, to get a token
}

// shared state for configuration
var conf *Config

// New returns the loaded configuration or panic
func New() (*Config, error) {
	var err error
	if conf == nil {
		conf, err = configInit("config.yaml")
	}
	return conf, err
}

// PROJECTNAME TODO: should be replaced in your application
const PROJECTNAME string = "go-cli"

// FIXME: not windows compatible
func configInit(filename string) (*Config, error) {
	viper := viper.New()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s", PROJECTNAME))
	viper.AddConfigPath(fmt.Sprintf("%s/.config/%s", os.ExpandEnv("$HOME"), PROJECTNAME))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("configuration format is not correct, caused by: %s", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	return &config, err
}
