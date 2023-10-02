package util

import (
	"net/http"
	"net/url"
)

func ClientWithProxy(proxyFunc func(reqURL *url.URL) (*url.URL, error), destUrl *url.URL) (*http.Client, error) {
	proxyUrl, err := proxyFunc(destUrl)
	if err != nil {
		return nil, err
	}
	proxy := http.ProxyURL(proxyUrl)
	return &http.Client{Transport: &http.Transport{Proxy: proxy}}, nil
}
