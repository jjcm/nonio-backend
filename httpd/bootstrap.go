package main

import (
	bs "github.com/jjcm/soci-backend/bootstrap"
	"github.com/jjcm/soci-backend/httpd/handlers"
	"github.com/jjcm/soci-backend/httpd/middleware"
	"github.com/jjcm/soci-backend/models"
)

var sociConfig bs.Config

func bootstrap() {
	c, err := bs.InitConfig()
	sociConfig = c
	if err != nil {
		logError(err)
		log(sociConfig)
		panic("Application can't start without a valid DB connection")
	}

	log("Application bootstrapped with these settings:")
	log("Port: " + sociConfig.AppPort)
	log("Database: " + sociConfig.DBDatabase)
	log("DB Username: " + sociConfig.DBUsername)

	// let's now hydrate a few things in the handlers package
	handlers.DBConn = sociConfig.DBConn
	handlers.Log = sociConfig.Logger
	handlers.HmacSampleSecret = sociConfig.HMACKey

	// let's now hydrate a few things in the middleware package
	middleware.Log = sociConfig.Logger

	// let's now hydrate a few things in the models package
	models.DBConn = sociConfig.DBConn
	models.Log = sociConfig.Logger
}
