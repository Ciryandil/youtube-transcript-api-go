package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

func loadCookieJar(filePath string) (*cookiejar.Jar, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	var cookieJson map[string][]*http.Cookie
	err = json.Unmarshal(data, &cookieJson)
	if err != nil {
		return nil, fmt.Errorf("error loading cookies from file data as json: %v", err)
	}
	var jar cookiejar.Jar
	for domain, cookies := range cookieJson {
		u := &url.URL{
			Scheme: "https",
			Host:   domain,
		}
		jar.SetCookies(u, cookies)
	}
	return &jar, nil
}
