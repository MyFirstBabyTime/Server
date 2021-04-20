package config

import (
	"github.com/spf13/viper"
	"time"
)

// App is the application config about syscheck domain
var App *authConfig

// init function initialize App global variable
func init() {
	App = &authConfig{}
}

// authConfig have config value and implement various interface about auth config
type authConfig struct {
	accessTokenDuration *time.Duration
}

// default const value about authConfig field
const (
	defaultAccessTokenDuration = time.Hour * 24
)

// AccessTokenDuration return access token valid duration
func (ac *authConfig) AccessTokenDuration() time.Duration {
	var key = "auth.accessTokenDuration"
	if ac.accessTokenDuration != nil {
		return *ac.accessTokenDuration
	}

	d, err := time.ParseDuration(viper.GetString(key))
	if err != nil {
		viper.Set(key, defaultAccessTokenDuration.String())
		d = defaultAccessTokenDuration
	}

	ac.accessTokenDuration = &d
	return *ac.accessTokenDuration
}
