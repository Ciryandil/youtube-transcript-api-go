package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func HTTPRequest(
	client *http.Client,
	method string,
	endpoint string,
	headers map[string]interface{},
	queryParams map[string]interface{},
	body io.Reader,
) (responseBody []byte, statusCode int, responseHeaders http.Header, err error) {

	parsedUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, 0, nil, err
	}
	query := parsedUrl.Query()
	for k, v := range queryParams {
		if v != nil {
			query.Set(k, fmt.Sprintf("%v", v))
		}
	}
	parsedUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(method, parsedUrl.String(), body)
	if err != nil {
		return nil, 0, nil, err
	}

	for k, v := range headers {
		if v != nil {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, resp.Header, err
	}

	return respBody, resp.StatusCode, resp.Header, err
}
