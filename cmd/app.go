package main

import (
	"fmt"
	"net/http"
	"time"

	"soci-backend/httpd"
	"soci-backend/httpd/middleware"
	"soci-backend/models"

	"github.com/go-co-op/gocron"
	"github.com/urfave/cli"
)

func runApp(c *cli.Context) error {
	for path, handler := range httpd.OpenRoutes() {
		http.HandleFunc(path, middleware.OpenCors(handler))
	}
	for path, handler := range httpd.ProtectedRoutes() {
		http.HandleFunc(path, middleware.OpenCors(middleware.CheckToken(handler)))
	}

	schedule := gocron.NewScheduler(time.UTC)
	//schedule.Every(1).Month(1).Do(models.AllocatePayouts)
	schedule.Every(1).Hour().Do(models.AllocatePayouts)

	schedule.StartAsync()
	_, nextTime := schedule.NextRun()
	log(fmt.Sprintf("Next payment calculation will run on %v", nextTime.Format("Mon Jan 2 15:04:05 2006")))

	log("Starting web api at port " + sociConfig.AppPort)
	http.ListenAndServe(":"+sociConfig.AppPort, nil)

	return nil
}
