// Package https provides a simple launcher for a Mux as HTTPS server.
// The certificate is retrieved from Let's Encrypt https://letsencrypt.org/
package https

import (
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func MustServeHttps(hostname string, mux *http.ServeMux) {
	mustStartRedirectHttp(mustStartHttps(hostname, mux))
}

func mustStartRedirectHttp(certManager *autocert.Manager) {
	redirectMux := &http.ServeMux{}
	redirectMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		newURI := "https://" + r.Host + r.URL.String()
		http.Redirect(w, r, newURI, http.StatusMovedPermanently)
	})
	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler: certManager.HTTPHandler(redirectMux),
	}
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}

func mustStartHttps(hostname string, mux *http.ServeMux) *autocert.Manager {
	certManager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			if host != hostname {
				return fmt.Errorf("acme/autocert: only %s host is allowed", hostname)
			}
			return nil
		},
		Cache: autocert.DirCache("."),
	}
	httpsSrv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
		Addr:         ":443",
		TLSConfig:    &tls.Config{GetCertificate: certManager.GetCertificate},
	}
	go func() {
		err := httpsSrv.ListenAndServeTLS("", "")
		if err != nil {
			panic(err)
		}
	}()
	return certManager
}

func MustServeHttp(mux *http.ServeMux) {
	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler: mux,
	}
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}

func WaitForIntOrTerm() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
