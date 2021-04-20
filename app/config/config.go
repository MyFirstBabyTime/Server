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

	// configFile represent aligo api key
	aligoAPIKey *string

	// configFile represent aligo account ID
	aligoAccountID *string

	// configFile represent aligo sender
	aligoSender *string

	// jwtKey represent jwt key
	jwtKey *string
}

// ConfigFile return config file get from environment variable
func (ac *appConfig) ConfigFile() string {
	if ac.configFile != nil {
		return *ac.configFile
	}

	if viper.IsSet("FIRST_BABY_TIME_CONFIG_FILE") {
		ac.configFile = _string(viper.GetString("FIRST_BABY_TIME_CONFIG_FILE"))
	} else {
		log.Fatal("please set FIRST_BABY_TIME_CONFIG_FILE in environment variable")
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

// AligoAPIKey return aligo api key get from environment variable
func (ac *appConfig) AligoAPIKey() string {
	if ac.aligoAPIKey != nil {
		return *ac.aligoAPIKey
	}

	if viper.IsSet("ALIGO_API_KEY") {
		ac.aligoAPIKey = _string(viper.GetString("ALIGO_API_KEY"))
	} else {
		log.Fatal("please set ALIGO_API_KEY in environment variable")
	}
	return *ac.aligoAPIKey
}

// AligoAccountID return aligo account ID key get from environment variable
func (ac *appConfig) AligoAccountID() string {
	if ac.aligoAccountID != nil {
		return *ac.aligoAccountID
	}

	if viper.IsSet("ALIGO_ACCOUNT_ID") {
		ac.aligoAccountID = _string(viper.GetString("ALIGO_ACCOUNT_ID"))
	} else {
		log.Fatal("please set ALIGO_ACCOUNT_ID in environment variable")
	}
	return *ac.aligoAccountID
}

// AligoSender return aligo sender  get from environment variable
func (ac *appConfig) AligoSender() string {
	if ac.aligoSender != nil {
		return *ac.aligoSender
	}

	if viper.IsSet("ALIGO_SENDER") {
		ac.aligoSender = _string(viper.GetString("ALIGO_SENDER"))
	} else {
		log.Fatal("please set ALIGO_SENDER in environment variable")
	}
	return *ac.aligoSender
}
// JwtKey return jwt key get from environment variable
func (ac *appConfig) JwtKey() string {
	if ac.jwtKey != nil {
		return *ac.jwtKey
	}

	if viper.IsSet("JWT_KEY") {
		ac.jwtKey = _string(viper.GetString("JWT_KEY"))
	} else {
		log.Fatal("please set JWT_KEY in environment variable")
	}
	return *ac.jwtKey
}

func _string(s string) *string { return &s }
