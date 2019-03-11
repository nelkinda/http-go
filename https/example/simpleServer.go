package main

import (
	"flag"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/details/uptime"
	"github.com/nelkinda/http-go/https"
	"net/http"
	"os"
)

func main() {
	serverNamePtr := flag.String("servername", os.Getenv("HOSTNMAE"), "Hostname for HTTPS.")
	startHttpsPtr := flag.Bool("https", false, "Start HTTPS.")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health.New(health.Health{Version: "1", ReleaseId: "0.0.1-SNAPSHOT"}, uptime.Process()).Handler)

	if *startHttpsPtr {
		https.MustServeHttps(mux, *serverNamePtr)
	} else {
		https.MustServeHttp(mux)
	}

	https.WaitForIntOrTerm()
}
