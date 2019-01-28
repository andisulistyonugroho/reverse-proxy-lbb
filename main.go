package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

/*
	Structs
*/

type requestPayloadStruct struct {
	ProxyCondition string `json:"proxy_condition"`
}

/*
	Entry
*/

func main() {
	// Log setup values
	logSetup()

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}

/*
	Utilities
*/

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {

	log.Printf("[%s]: %s FROM %s\n", req.Method, req.URL, req.RemoteAddr)

	if req.Method != "POST" {
		req.ParseForm()
		log.Printf("PARAM: %s\n", req.Form)
	}
	log.Printf("======================================================================\n")

	url := "https://api.ligabukubandung.com"
	serveReverseProxy(url, res, req)
}

/*
	Reverse Proxy Logic
*/

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

/*
	Getters
*/

// Get the port to listen on
func getListenAddress() string {
	port := getEnv("PORT", "1338")
	return ":" + port
}

/*
	Logging
*/

// Log the env variables required for a reverse proxy
func logSetup() {
	log.Printf("Server will run on: %s\n", getListenAddress())
	log.Printf("======================================================================\n")
}
