# `http-go` - Making http even more useful

## HTTPS
package `https` provides a function to serve HTTPS with Let's Encrypt.
All you need is this package and a mapped hostname, and you're ready to go with just one function call.

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

## MIME Types
package `mimetypes` contains all MIME Types registered with IANA.
