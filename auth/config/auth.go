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
	// fields using in auth domain usecase (implement authUsecaseConfig)
	// accessTokenDuration represent time valid duration for access token
	accessTokenDuration *time.Duration

	// parentProfileS3Bucket represent aws s3 bucket for parent profile
	parentProfileS3Bucket *string
}

// default const value about authConfig field
const (
	defaultAccessTokenDuration   = time.Hour * 24
	defaultParentProfileS3Bucket = "first-baby-time"
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

// ParentProfileS3Bucket implement ParentProfileS3Bucket of authUsecaseConfig
func (ac *authConfig) ParentProfileS3Bucket() string {
	var key = "auth.parentProfileS3Bucket"
	if ac.parentProfileS3Bucket == nil {
		if _, ok := viper.Get(key).(string); !ok {
			viper.Set(key, defaultParentProfileS3Bucket)
		}
		ac.parentProfileS3Bucket = _string(viper.GetString(key))
	}
	return *ac.parentProfileS3Bucket
}

func _string(s string) *string { return &s }
