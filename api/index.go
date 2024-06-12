package handler

import (
	"io"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-PROXY-HOST, X-PROXY-SCHEME")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	loadProxyParams := func(r *http.Request, key string) string {
		value := r.Header.Get(key)
		if value == "" {
			value = r.URL.Query().Get(key)
			r.URL.Query().Del(key)
		}
		r.Header.Del(key)
		return value
	}

	hostname := loadProxyParams(r, "X-PROXY-HOST")
	proxyScheme := loadProxyParams(r, "X-PROXY-SCHEME")
	if hostname == "" {
		http.Error(w, "X-PROXY-HOST is required", http.StatusBadRequest)
		return
	}
	if proxyScheme == "" {
		proxyScheme = "https"
	}
	r.URL.Host = hostname
	r.URL.Scheme = proxyScheme

	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
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
