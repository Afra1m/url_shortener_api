package utils

import (
	"net/url"
)

func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func GenerateBaseURL(host, port string) string {
	if port == "" || port == "80" {
		return "http://" + host
	}
	return "http://" + host + ":" + port
}
