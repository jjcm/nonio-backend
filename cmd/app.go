package main

import (
	"net/http"

	"soci-backend/httpd"
	"soci-backend/httpd/middleware"

	"github.com/urfave/cli"
)

func runApp(c *cli.Context) error {
	for path, handler := range httpd.OpenRoutes() {
		http.HandleFunc(path, handler)
	}
	for path, handler := range httpd.ProtectedRoutes() {
		http.HandleFunc(path, middleware.CheckToken(middleware.OpenCors(handler)))
	}
	log("Starting web api at port " + sociConfig.AppPort)
	http.ListenAndServe(":"+sociConfig.AppPort, nil)

	return nil
}
