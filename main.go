package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getTargetUrl(orgUrl *url.URL) (*url.URL, error) {
	path := strings.TrimPrefix(orgUrl.Path, "/")
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		path = "https://" + path
	}
	if orgUrl.RawQuery != "" {
		path += "?" + orgUrl.RawQuery
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	if u.Host == "" || !strings.Contains(u.Host, ".") {
		return nil, fmt.Errorf("invalid host")
	}

	return u, nil
}

type OneHandler struct {
}

func (h *OneHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	targetUrl, err := getTargetUrl(r.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fwReq := &http.Request{
		Method: r.Method,
		URL:    targetUrl,
		Body:   r.Body,
		Header: r.Header,
	}

	log.Println(r.Method, targetUrl.String())

	fwReq.Header["Host"] = []string{targetUrl.Host}

	client := &http.Client{}
	resp, err := client.Do(fwReq)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, &OneHandler{}))
}
