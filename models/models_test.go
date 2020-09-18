package models

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	bs "soci-backend/bootstrap"

	_ "github.com/go-sql-driver/mysql"
)

func setupTestingDB() error {
	var testingDBName = "socidb_testing"
	os.Setenv("APP_KEY", "secret")
	os.Setenv("OAUTH_ID", "12345")
	os.Setenv("OAUTH_SECRET", "12345")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "dbtestuser")
	os.Setenv("DB_DATABASE", testingDBName)
	os.Setenv("DB_PASSWORD", "password")

	c, err := bs.InitConfig()
	if err != nil {
		fmt.Println("bootstrap died", err)
		return err
	}

	os.Setenv("DB_DATABASE", testingDBName)
	c, err = bs.InitConfig()
	if err != nil {
		fmt.Println("db connection died")
		panic(err)
	}
	DBConn = c.DBConn
	Log = c.Logger // so we don't choke on any log calls

	// get the database back to square 1
	resetTestingDB()

	goPath := os.Getenv("GOPATH")
	command := goPath + "/bin/goose"
	cmd := exec.Command(command, "mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@/"+testingDBName, "up")

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	workingDir = strings.Replace(workingDir, "models", "migrations", -1)
	cmd.Dir = workingDir

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
