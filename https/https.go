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

func MustServeHttps(certDir string, mux *http.ServeMux, hostnames ...string) {
	mustStartRedirectHttp(mustStartHttps(certDir, mux, hostnames...))
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

func mustStartHttps(certDir string, mux *http.ServeMux, hostnames ...string) *autocert.Manager {
	certManager := &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: func(ctx context.Context, host string) error {
			for _, h := range hostnames {
				if h == host {
					return nil
				}
			}
			return fmt.Errorf("acme/autocert: %s not in the list of allowed hostnames %v", host, hostnames)
		},
		Cache: autocert.DirCache(certDir),
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

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func LogHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}
		handler.ServeHTTP(&sw, r)
		duration := time.Now().Sub(start)
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s: served %s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%s\t\"%s\"\t%d\n", os.Args[0], "info", start.UTC().Format(time.RFC3339), r.Host, r.RemoteAddr, r.Method, r.RequestURI, r.Proto, sw.status, sw.length, r.Referer(), r.Header.Get("User-Agent"), duration)
	}
}

func LogHandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	return LogHandler(handler)
}
