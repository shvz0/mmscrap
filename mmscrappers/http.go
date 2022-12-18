package mmscrappers

import (
	"fmt"
	"io"
	"net/http"
)

func GetPage(url string, headers map[string][]string) (io.Reader, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header = headers

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	return resp.Body, err
}

func defaultHeaders() map[string][]string {
	return map[string][]string{
		"User-Agent":                {"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		"Accept-Language":           {"en-GB,en;q=0.5"},
		"Connection":                {"keep-alive"},
		"Upgrade-Insecure-Requests": {"1"},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"cross-site"},
	}
}
