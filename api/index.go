package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func internalServerError(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("Internal server error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("WithHandler panic: %v", err)
			http.Error(w, fmt.Sprintf("internal server error: %v", err), http.StatusInternalServerError)
		}
	}()

	htmlProxy := os.Getenv("HTTP_PROXY_ENABLE") == "true"

	// Set the CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-PROXY-HOST, X-PROXY-SCHEME")

	// Handle the OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Redirect to the GitHub repository
	if r.URL.Path == "/" {
		http.Redirect(w, r, "https://github.com/TBXark/vercel-proxy", http.StatusMovedPermanently)
		return
	}

	// Get the URL to proxy
	re := regexp.MustCompile(`^/*(https?:)/*`)
	u := re.ReplaceAllString(r.URL.Path, "$1//")
	if r.URL.RawQuery != "" {
		u += "?" + r.URL.RawQuery
	}
	if !strings.HasPrefix(u, "http") {
		http.Error(w, "invalid url: "+u, http.StatusBadRequest)
		return
	}

	// Create a new request
	req, err := http.NewRequest(r.Method, u, r.Body)
	if err != nil {
		internalServerError(w, err)
		return
	}
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	if htmlProxy && r.Header.Get("Accept-Encoding") != "" {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			req.Header.Set("Accept-Encoding", "gzip")
		}
	}

	// Send the request to the real server
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		internalServerError(w, err)
		return
	}
	defer func(writer http.ResponseWriter, response *http.Response) {
		internalServerError(writer, response.Body.Close())
	}(w, resp)

	if e := proxyRaw(w, resp, r); e != nil {
		internalServerError(w, e)
		return
	}

	w.WriteHeader(resp.StatusCode)
}

func proxyRaw(w http.ResponseWriter, resp *http.Response, req *http.Request) error {
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	if w.Header().Get("Referer") != "" {
		w.Header().Del("Referer")
		w.Header().Add("Referer", req.Host)
	}

	// Copy the response body to the output stream
	_, err := io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
