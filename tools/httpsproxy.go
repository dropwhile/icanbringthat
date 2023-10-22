package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	localProxyUrl, _ := url.Parse("http://127.0.0.1:8000/")
	localProxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(localProxyUrl)
			r.SetXForwarded()
		},
	}
	//localProxy := httputil.NewSingleHostReverseProxy(localProxyUrl)
	http.Handle("/", localProxy)

	log.Println("Serving on localhost:8080")
	log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
}
