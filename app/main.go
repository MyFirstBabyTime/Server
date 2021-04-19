package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"

	"github.com/MyFirstBabyTime/Server/hash"
	"github.com/MyFirstBabyTime/Server/jwt"
	"github.com/MyFirstBabyTime/Server/message"
	"github.com/MyFirstBabyTime/Server/parser"
	"github.com/MyFirstBabyTime/Server/tx"
	"github.com/MyFirstBabyTime/Server/validate"

	_authHttpDelivery "github.com/MyFirstBabyTime/Server/auth/delivery/http"
	_authRepo "github.com/MyFirstBabyTime/Server/auth/repository/mysql"
	_authUcase "github.com/MyFirstBabyTime/Server/auth/usecase"
)

func init() {
	// set flag to log current date, time & long file name
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "", "", "", "", "")
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create mysql connection").Error())
	}

	r := gin.Default()

	_ps := parser.MysqlMsgParser()
	_vl := validate.New()
	_tx := tx.NewSqlxHandler(db)
	_msg := message.AligoAgent("", "", "")
	_hash := hash.BcryptHandler()
	_jwt := jwt.UUIDHandler("")

	au := _authUcase.AuthUsecase(
		_authRepo.ParentAuthRepository(db, _ps, _vl),
		_authRepo.ParentPhoneCertifyRepository(db, _ps, _vl),
		_tx, _msg, _hash, _jwt,
	)
	_authHttpDelivery.NewAuthHandler(r, au, _vl)

	log.Fatal(r.Run(":8000"))
}
