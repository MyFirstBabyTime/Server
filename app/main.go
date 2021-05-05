package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"

	"github.com/MyFirstBabyTime/Server/app/config"
	"github.com/MyFirstBabyTime/Server/hash"
	"github.com/MyFirstBabyTime/Server/jwt"
	"github.com/MyFirstBabyTime/Server/message"
	"github.com/MyFirstBabyTime/Server/parser"
	"github.com/MyFirstBabyTime/Server/s3"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/MyFirstBabyTime/Server/validate"

	_authConfig "github.com/MyFirstBabyTime/Server/auth/config"
	_authHttpDelivery "github.com/MyFirstBabyTime/Server/auth/delivery/http"
	_authRepo "github.com/MyFirstBabyTime/Server/auth/repository/mysql"
	_authUcase "github.com/MyFirstBabyTime/Server/auth/usecase"

	_expenditureDelivery "github.com/MyFirstBabyTime/Server/chlidcare-expenditure/delivery/http"
	_expenditureRepo "github.com/MyFirstBabyTime/Server/chlidcare-expenditure/repository/mysql"
	_expenditureUcase "github.com/MyFirstBabyTime/Server/chlidcare-expenditure/usecase"

	_cloudMaintainerDelivery "github.com/MyFirstBabyTime/Server/cloud-maintainer/delivery/http"
	_cloudMaintainerUsecase "github.com/MyFirstBabyTime/Server/cloud-maintainer/usecase"
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

	s3Ses, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.App.S3Region()),
		Credentials: credentials.NewStaticCredentials(config.App.AwsS3ID(), config.App.AwsS3Key(), ""),
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create aws connection").Error())
	}

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization", "authorization", "Request-Security")

	r.Use(cors.New(corsConfig))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	_ps := parser.MysqlMsgParser()
	_vl := validate.New()
	_tx := tx.NewSqlxHandler(db)
	_msg := message.AligoAgent(config.App.AligoAPIKey(), config.App.AligoAccountID(), config.App.AligoSender())
	_hash := hash.BcryptHandler()
	_jwt := jwt.UUIDHandler(config.App.JwtKey())
	_s3 := s3.New(s3Ses)

	au := _authUcase.AuthUsecase(
		_authConfig.App,
		_authRepo.ParentAuthRepository(_authConfig.App, db, _ps, _vl),
		_authRepo.ParentPhoneCertifyRepository(_authConfig.App, db, _ps, _vl),
		_tx, _msg, _hash, _jwt, _s3,
	)
	_authHttpDelivery.NewAuthHandler(r, au, _vl)

	eu := _expenditureUcase.ExpenditureUsecase(
		_expenditureRepo.ExpenditureRepository(db, _ps, _vl),
		_tx,
	)
	_expenditureDelivery.NewExpenditureHandler(r, eu, _vl, _jwt,)

	cmu := _cloudMaintainerUsecase.CloudMaintainerUsecase(
		config.App.CloudManagementKey(),
	)
	_cloudMaintainerDelivery.NewCloudMaintainerHandler(r, cmu, _vl)

	log.Fatal(r.Run(":80"))
}
