package base

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

// AppName exposes application name to config module
var AppName string

// Config receives configuration options
type Config struct {
	mutex          sync.Mutex
	File           string
	Host           string  `yaml:"host"`
	Port           int     `yaml:"port"`
	User           string  `yaml:"user"`
	Password       string  `yaml:"password"`
	Database       string  `yaml:"database"`
	Interval       float64 `yaml:"interval"`
	ConnectTimeout int     `yaml:"connect-timeout"`
	IdleTimeout    float64 `yaml:"idle-timeout"`
	ActiveTimeout  float64 `yaml:"active-timeout"`
	LogFile        string  `yaml:"log-file"`
	PidFile        string  `yaml:"pid-file"`
}

func init() {
	AppName = "pgterminate"
}

// NewConfig creates a Config object
func NewConfig() *Config {
	return &Config{}
}

// Read loads options from a configuration file to Config
func (c *Config) Read(file string) error {
	file, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return err
	}

	return nil
}

// Reload reads from file and update configuration
func (c *Config) Reload() {
	log.Println("Reloading configuration")
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.File != "" {
		c.Read(c.File)
	}
}

// Dsn formats a connection string based on Config
func (c *Config) Dsn() string {
	var parameters []string
	if c.Host != "" {
		parameters = append(parameters, fmt.Sprintf("host=%s", c.Host))
	}
	if c.Port != 0 {
		parameters = append(parameters, fmt.Sprintf("port=%d", c.Port))
	}
	if c.User != "" {
		parameters = append(parameters, fmt.Sprintf("user=%s", c.User))
	}
	if c.Password != "" {
		parameters = append(parameters, fmt.Sprintf("password=%s", c.Password))
	}
	if c.Database != "" {
		parameters = append(parameters, fmt.Sprintf("database=%s", c.Database))
	}
	if c.ConnectTimeout != 0 {
		parameters = append(parameters, fmt.Sprintf("connect_timeout=%d", c.ConnectTimeout))
	}
	if AppName != "" {
		parameters = append(parameters, fmt.Sprintf("application_name=%s", AppName))
	}
	return strings.Join(parameters, " ")
}
