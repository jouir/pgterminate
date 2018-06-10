package base

// Context stores dynamic values like channels and exposes configuration
type Context struct {
	Sessions chan Session
	Done     chan bool
	Config   *Config
}

// NewContext instanciates a Context
func NewContext(config *Config, sessions chan Session, done chan bool) *Context {
	return &Context{
		Config:   config,
		Sessions: sessions,
		Done:     done,
	}
}
