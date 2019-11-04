package models

import (
	"bytes"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	bs "github.com/jjcm/soci-backend/bootstrap"
)

func setupTestingDB() error {
	var testingDBName = "socidb_testing"
	os.Setenv("APP_KEY", "secret")
	os.Setenv("OAUTH_ID", "12345")
	os.Setenv("OAUTH_SECRET", "12345")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "genius")

	c, err := bs.InitConfig()
	if err != nil {
		return err
	}

	tempDB := c.DBConn
	tempDB.Exec("CREATE DATABASE " + testingDBName)
	tempDB.Close()

	os.Setenv("DB_DATABASE", testingDBName)
	c, err = bs.InitConfig()
	DBConn = c.DBConn
	if err != nil {
		panic(err)
	}
	Log = c.Logger // so we don't choke on any log calls

	cmd := exec.Command("/home/lapubell/programming/go/bin/goose", "mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@/"+testingDBName, "up")
	cmd.Dir = "/home/lapubell/programming/go/src/github.com/jjcm/soci-backend/migrations"
	var output bytes.Buffer
	cmd.Stderr = &output
	cmd.Run()

	return nil
}

func teardownTestingDB() {
	var tables []string
	DBConn.Select(&tables, "SHOW TABLES")
	for _, t := range tables {
		DBConn.Exec("DROP TABLE " + t)
	}
}
