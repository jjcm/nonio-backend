package main

import (
	"net/http"

	"github.com/jjcm/soci-backend/httpd/middleware"
	"github.com/urfave/cli"
)

func runApp(c *cli.Context) error {
	for path, handler := range openRoutes() {
		http.HandleFunc(path, handler)
	}
	for path, handler := range protectedRoutes() {
		http.HandleFunc(path, middleware.CheckToken(handler))
	}
	log("Starting web api at port " + sociConfig.AppPort)
	http.ListenAndServe(":"+sociConfig.AppPort, nil)

	return nil
}
