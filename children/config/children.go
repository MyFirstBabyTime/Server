package config

import "github.com/spf13/viper"

// App is the application config about children domain
var App *childrenConfig

// init function initialize App global variable
func init() {
	App = &childrenConfig{}
}

// childrenConfig have config value and implement various interface about children config
type childrenConfig struct {
	// childrenProfileS3Bucket represent aws s3 bucket for chlidren profile
	childrenProfileS3Bucket *string
}

// default const value about childrenConfig field
const (
	defaultChildrenProfileS3Bucket = "first-baby-time"
)

// ChildrenProfileS3Bucket implement ChildrenProfileS3Bucket of childrenUsecaseConfig
func (cc *childrenConfig) ChildrenProfileS3Bucket() string {
	var key = "children.childrenProfileS3Bucket"
	if cc.childrenProfileS3Bucket == nil {
		if _, ok := viper.Get(key).(string); !ok {
			viper.Set(key, defaultChildrenProfileS3Bucket)
		}
		cc.childrenProfileS3Bucket = _string(viper.GetString(key))
	}
	return *cc.childrenProfileS3Bucket
}

func _string(s string) *string { return &s }
