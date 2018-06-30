package base

import (
	"fmt"
	"github.com/jouir/pgterminate/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// AppName exposes application name to config module
var AppName string

// Config receives configuration options
type Config struct {
	mutex                     sync.Mutex
	File                      string
	Host                      string      `yaml:"host"`
	Port                      int         `yaml:"port"`
	User                      string      `yaml:"user"`
	Password                  string      `yaml:"password"`
	Database                  string      `yaml:"database"`
	Interval                  float64     `yaml:"interval"`
	ConnectTimeout            int         `yaml:"connect-timeout"`
	IdleTimeout               float64     `yaml:"idle-timeout"`
	ActiveTimeout             float64     `yaml:"active-timeout"`
	LogDestination            string      `yaml:"log-destination"`
	LogFile                   string      `yaml:"log-file"`
	PidFile                   string      `yaml:"pid-file"`
	SyslogIdent               string      `yaml:"syslog-ident"`
	SyslogFacility            string      `yaml:"syslog-facility"`
	IncludeUsers              StringFlags `yaml:"include-users"`
	IncludeUsersRegex         string      `yaml:"include-users-regex"`
	IncludeUsersRegexCompiled *regexp.Regexp
	ExcludeUsers              StringFlags `yaml:"exclude-users"`
	ExcludeUsersRegex         string      `yaml:"exclude-users-regex"`
	ExcludeUsersRegexCompiled *regexp.Regexp
	Cancel                    bool `yaml:"cancel"`
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

// Reload reads from file to update configuration and re-compile regexes
func (c *Config) Reload() {
	log.Debug("Reloading configuration")
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.File != "" {
		c.Read(c.File)
	}
	err := c.CompileRegexes()
	Panic(err)
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

// CompileRegexes transforms regexes from string to regexp instance
func (c *Config) CompileRegexes() (err error) {
	if c.IncludeUsersRegex != "" {
		c.IncludeUsersRegexCompiled, err = regexp.Compile(c.IncludeUsersRegex)
		if err != nil {
			return err
		}
	}
	if c.ExcludeUsersRegex != "" {
		c.ExcludeUsersRegexCompiled, err = regexp.Compile(c.ExcludeUsersRegex)
		if err != nil {
			return err
		}
	}
	return nil
}

// StringFlags append multiple string flags into a string slice
type StringFlags []string

// String for implementing flag interface
func (s *StringFlags) String() string {
	return "multiple strings flag"
}

// Set adds alues into the slice
func (s *StringFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}
