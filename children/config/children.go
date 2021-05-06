package config

// App is the application config about syscheck domain
var App *authConfig

// init function initialize App global variable
func init() {
	App = &authConfig{}
}

// authConfig have config value and implement various interface about auth config
type authConfig struct {}

func _string(s string) *string { return &s }
