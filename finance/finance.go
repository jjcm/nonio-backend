package finance

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// DBConn this is here so that we can hydrate all the web handlers with the same
// DB connection that we are using in the main package
var DBConn *sqlx.DB

var Log *logrus.Logger
var ServerFee float64
