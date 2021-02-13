package main

import (
	"fmt"
	"net/http"
	"time"

	"soci-backend/finance"
	"soci-backend/httpd"
	"soci-backend/httpd/middleware"

	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli"
)

func runApp(c *cli.Context) error {
	for path, handler := range httpd.OpenRoutes() {
		http.HandleFunc(path, middleware.OpenCors(handler))
	}
	for path, handler := range httpd.ProtectedRoutes() {
		http.HandleFunc(path, middleware.OpenCors(middleware.CheckToken(handler)))
	}

	t := time.Date(2020, time.January, 1, 0, 28, 0, 0, time.Local)
	gocron.Every(12).Weeks().From(&t).Do(finance.CalculatePayouts)

	_, nextTime := gocron.NextRun()
	fmt.Println("next run is at...")
	fmt.Println(nextTime)
	fmt.Println(":)")
	gocron.Start()

	log("Starting web api at port " + sociConfig.AppPort)
	http.ListenAndServe(":"+sociConfig.AppPort, nil)

	return nil
}
