package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	listen := flag.String("listen", ":8080", "address to listen on")
	target := flag.String("target", "", "upstream URL (for example https://api.example.com)")
	origin := flag.String("origin", "http://localhost:5500", "allowed frontend origin")
	flag.Parse()
	if *target == "" {
		log.Fatal("-target is required")
	}
	upstream, err := url.Parse(*target)
	if err != nil || upstream.Scheme == "" || upstream.Host == "" {
		log.Fatalf("invalid upstream URL: %q", *target)
	}

	proxy := httputil.NewSingleHostReverseProxy(upstream)
	director := proxy.Director
	proxy.Director = func(r *http.Request) {
		director(r)
		r.Host = upstream.Host
	}
	proxy.ErrorLog = log.New(os.Stderr, "proxy: ", log.LstdFlags)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") == *origin {
			w.Header().Set("Access-Control-Allow-Origin", *origin)
			w.Header().Add("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		proxy.ServeHTTP(w, r)
	})
	log.Printf("CORS development proxy listening on %s -> %s", *listen, upstream)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
