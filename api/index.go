package api

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-PROXY-HOST, X-PROXY-SCHEME")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	re := regexp.MustCompile(`^/{0,}(https?:)/{0,}`)
	u := re.ReplaceAllString(r.URL.Path, "$1//")
	if r.URL.RawQuery != "" {
		u += "?" + r.URL.RawQuery
	}

	if !strings.HasPrefix(u, "http") {
		http.Error(w, "invalid url: "+u, http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, u, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
}
