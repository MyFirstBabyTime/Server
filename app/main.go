package main

import (
	"github.com/MyFirstBabyTime/Server/app/config"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"

	"github.com/MyFirstBabyTime/Server/hash"
	"github.com/MyFirstBabyTime/Server/jwt"
	"github.com/MyFirstBabyTime/Server/message"
	"github.com/MyFirstBabyTime/Server/parser"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/MyFirstBabyTime/Server/validate"

	_authConfig "github.com/MyFirstBabyTime/Server/auth/config"
	_authHttpDelivery "github.com/MyFirstBabyTime/Server/auth/delivery/http"
	_authRepo "github.com/MyFirstBabyTime/Server/auth/repository/mysql"
	_authUcase "github.com/MyFirstBabyTime/Server/auth/usecase"
)

func init() {
	// set flag to log current date, time & long file name
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// set and read config file in viper package
	viper.AutomaticEnv()
	viper.SetConfigFile(config.App.ConfigFile())
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	db, err := sqlx.Connect("mysql", config.App.MysqlDataSource())
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create mysql connection").Error())
	}

	r := gin.Default()

	_ps := parser.MysqlMsgParser()
	_vl := validate.New()
	_tx := tx.NewSqlxHandler(db)
	_msg := message.AligoAgent(config.App.AligoAPIKey(), config.App.AligoAccountID(), config.App.AligoSender())
	_hash := hash.BcryptHandler()
	_jwt := jwt.UUIDHandler(config.App.JwtKey())

	au := _authUcase.AuthUsecase(
		_authConfig.App,
		_authRepo.ParentAuthRepository(_authConfig.App, db, _ps, _vl),
		_authRepo.ParentPhoneCertifyRepository(_authConfig.App, db, _ps, _vl),
		_tx, _msg, _hash, _jwt,
	)
	_authHttpDelivery.NewAuthHandler(r, au, _vl)

	log.Fatal(r.Run(":8000"))
}
