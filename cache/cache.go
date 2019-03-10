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
	Gzip = "gzip"
)

type Entry struct {
	URI          string
	Body         []byte
	GzipBody     []byte
	ContentType  string
	LastModified *time.Time
	MaxAge       time.Duration
	ETag         string
}

type Cache struct {
	Cache map[string]*Entry
}

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
	modTime := fileStat.ModTime()
	c.Cache[uri] = &Entry{
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
