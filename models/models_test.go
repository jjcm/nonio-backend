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

	os.Setenv("DB_DATABASE", testingDBName)
	c, err = bs.InitConfig()
	if err != nil {
		panic(err)
	}
	DBConn = c.DBConn
	Log = c.Logger // so we don't choke on any log calls

	// get the database back to square 1
	resetTestingDB()

	cmd := exec.Command("/home/lapubell/programming/go/bin/goose", "mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@/"+testingDBName, "up")
	cmd.Dir = "/home/lapubell/programming/go/src/github.com/jjcm/soci-backend/migrations"
	var output bytes.Buffer
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		panic(output.String())
	}

	return nil
}

func resetTestingDB() {
	DBConn.Exec("SET FOREIGN_KEY_CHECKS=0")

	var tables []string
	DBConn.Select(&tables, "SHOW TABLES")
	for _, t := range tables {
		_, err := DBConn.Exec("DROP TABLE " + t)
		if err != nil {
			panic(err)
		}
	}
	DBConn.Exec("SET FOREIGN_KEY_CHECKS=1")
}
