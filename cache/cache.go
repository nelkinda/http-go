package cache

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/nelkinda/http-go/header"
	"github.com/nelkinda/http-go/mimetype"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// Gzip is the constant for the Content-Encoding "gzip".
	Gzip = "gzip"
)

// Entry describes a single entry in the cache.
type Entry struct {

	// The URI to which this entry is mapped.
	// This is the same as the key in the cache.
	URI          string

	// The uncompressed response body.
	Body         []byte

	// The GZip compressed response body.
	GzipBody     []byte

	// The Content-Type of the body.
	ContentType  string

	// The Last-Modified timestamp of the body.
	LastModified *time.Time

	// The maximum cache age of the body.
	MaxAge       time.Duration

	// The ETag of the body, if available.
	ETag         string
}

// Cache is a HTTP Cache.
type Cache struct {
	Cache map[string]*Entry
}

// GlobalCache is the global (default) cache.
var GlobalCache = &Cache{Cache: make(map[string]*Entry)}

func CacheHandlerFunc(fallback http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		relativePath := r.RequestURI[1:len(r.URL.Path)]
		switch relativePath {
		case "":
			fallback.ServeHTTP(w, r)
		default:
			ServeCacheEntry(w, r, relativePath)
		}
	}
}

func CacheHandler(fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		relativePath := r.RequestURI[1:len(r.URL.Path)]
		switch relativePath {
		case "":
			fallback.ServeHTTP(w, r)
		default:
			ServeCacheEntry(w, r, relativePath)
		}
	}
}

func ServeCacheEntry(w http.ResponseWriter, r *http.Request, id string) {
	if cacheEntry, ok := GlobalCache.Cache[id]; !ok {
		http.NotFoundHandler().ServeHTTP(w, r)
	} else {
		cacheEntry.Serve(w, r)
	}
}

func (e *Entry) Serve(w http.ResponseWriter, r *http.Request) {
	contentType := e.ContentType
	contentType = fixContentType(r, contentType)
	w.Header().Add(header.ContentType, contentType)
	if e.ETag == "" {
		i := md5.Sum(e.Body)
		e.ETag = `"` + hex.EncodeToString(i[:]) + `"`
	}
	w.Header().Add(header.ETag, e.ETag)
	if e.LastModified != nil {
		w.Header().Add(header.LastModified, e.LastModified.Format(http.TimeFormat))
	}
	if e.MaxAge != 0 {
		w.Header().Add(header.Expires, time.Now().Add(e.MaxAge).Format(http.TimeFormat))
		w.Header().Add(header.CacheControl, fmt.Sprintf("max-age=%d", int(e.MaxAge.Seconds())))
	}
	if isGzip(r) {
		w.Header().Add(header.ContentEncoding, Gzip)
		_, _ = w.Write(e.GzipBody)
	} else {
		_, _ = w.Write(e.Body)
	}
}

func fixContentType(request *http.Request, s string) string {
	if s != mimetype.ApplicationXhtmlXml {
		return s
	}
	userAgent := request.Header.Get("User-Agent")
	if strings.Contains(userAgent, "Twitter") || strings.Contains(userAgent, "LinkedIn") {
		return "text/html"
	}
	return s
}

func isGzip(r *http.Request) bool {
	acceptEncoding := r.Header.Get(header.AcceptEncoding)
	encodings := strings.Split(acceptEncoding, ",")
	for _, encoding := range encodings {
		if strings.TrimSpace(encoding) == Gzip {
			return true
		}
	}
	return false
}

func (c *Cache) LoadCacheFile(filename string, uri string, contentType string, maxAge time.Duration) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fileStat, err := os.Stat(filename)
	if err != nil {
		return err
	}
	modTime := fileStat.ModTime().UTC()
	c.Cache[uri] = &Entry{
		URI:          uri,
		Body:         body,
		GzipBody:     compressGzip(body),
		ContentType:  contentType,
		LastModified: &modTime,
		MaxAge:       maxAge,
	}
	return nil
}

func LoadCacheFile(filename string, uri string, contentType string, maxAge time.Duration) error {
	return GlobalCache.LoadCacheFile(filename, uri, contentType, maxAge)
}

func compressGzip(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func (c *Cache) Size() (entries int, memory int) {
	for _, entry := range c.Cache {
		entries++
		memory += len(entry.Body) + len(entry.GzipBody)
	}
	return entries, memory
}

func Size() (entries int, memory int) {
	return GlobalCache.Size()
}

func (c *Cache) Add(entry *Entry) {
	if entry.GzipBody == nil {
		entry.GzipBody = compressGzip(entry.Body)
	}
	c.Cache[entry.URI] = entry
}

func Add(entry *Entry) {
	GlobalCache.Add(entry)
}

func (c *Cache) Sitemap(r *http.Request) string {
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">
`
	for key, entry := range c.Cache {
		switch entry.ContentType {
		case mimetype.ApplicationXhtmlXml, mimetype.TextHtml:
			if entry.LastModified != nil {
				sitemap += fmt.Sprintf("<url><loc>https://%s%s</loc><lastmod>%s</lastmod></url>\n", r.Host, key, entry.LastModified.Format(time.RFC3339Nano))
			} else {
				sitemap += fmt.Sprintf("<url><loc>https://%s%s</loc></url>\n", r.Host, entry.URI)
			}
		}
	}
	sitemap += `</urlset>`
	return sitemap
}

func Sitemap(r *http.Request) string {
	return GlobalCache.Sitemap(r)
}