package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// App is the application config using in main package
var App *appConfig

// init function initialize App global variable
func init() {
	App = &appConfig{}
}

// appConfig having config value and return that value with method. Not implement interface
type appConfig struct {
	// configFile represent full name of config file
	configFile *string

	// mysqlDataSource represent data source name of MySQL
	mysqlDataSource *string
}

// ConfigFile return config file get from environment variable
func (ac *appConfig) ConfigFile() string {
	if ac.configFile != nil {
		return *ac.configFile
	}

	if viper.IsSet("CONFIG_FILE") {
		ac.configFile = _string(viper.GetString("CONFIG_FILE"))
	} else {
		log.Fatal("please set CONFIG_FILE in environment variable")
	}
	return *ac.configFile
}


// MysqlDataSource return mysql data source name with value get from environment variable
func (ac *appConfig) MysqlDataSource() string {
	if ac.mysqlDataSource != nil {
		return *ac.mysqlDataSource
	}

	format := "%s:%s@tcp(%s)/%s"
	var args []interface{}

	if viper.IsSet("MYSQL_USERNAME") {
		args = append(args, viper.GetString("MYSQL_USERNAME"))
	} else {
		log.Fatal("please set MYSQL_USERNAME in environment variable")
	}

	if viper.IsSet("MYSQL_PASSWORD") {
		args = append(args, viper.GetString("MYSQL_PASSWORD"))
	} else {
		log.Fatal("please set MYSQL_PASSWORD in environment variable")
	}

	if viper.IsSet("MYSQL_ADDRESS") {
		args = append(args, viper.GetString("MYSQL_ADDRESS"))
	} else {
		log.Fatal("please set MYSQL_ADDRESS in environment variable")
	}

	if viper.IsSet("MYSQL_DATABASE") {
		args = append(args, viper.GetString("MYSQL_DATABASE"))
	} else {
		log.Fatal("please set MYSQL_DATABASE in environment variable")
	}

	ac.mysqlDataSource = _string(fmt.Sprintf(format, args...))
	return *ac.mysqlDataSource
}

func _string(s string) *string { return &s }
