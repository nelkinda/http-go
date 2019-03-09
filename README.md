# `http-go` - Making http even more useful

## HTTPS
package `https` provides a function to serve HTTPS with Let's Encrypt.
All you need is this package and a mapped hostname, and you're ready to go with just one function call.
The function `MustServeHttps` will
* Configure the Certificate Manager to obtain the certificate from Let's Encrypt.
* Start the HTTPS server in a goroutine based on the `http.ServeMux` provided as argument.
* Start an HTTP server in a goroutine which redirects to HTTPS as well as handles the Let's Encrypt callback.

Because the servers are started in goroutines, the main function now needs to wait for termination.
The `https` package provides a utility function for that.

Example code:
```go
package main

import (
	"github.com/nelkinda/http-go/https"
	"net/http"
	"os"
)

func main() {
	mux := createMux()
	https.MustServeHttps("myhost.com", mux)
	https.WaitForIntOrTerm()
	os.Exit(0)
}

func createMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("…", …)
	return mux
}
```

This will start an HTTPS server on `myhost.com`, requesting a certificate from Let's Encrypt at the start.

## Headers
package `headers` contains all HTTP headers defined in the HTTP specifications.

## MIME Types
package `mimetypes` contains all MIME Types registered with IANA.
