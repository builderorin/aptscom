package client

import (
	"io"
	"time"

	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/fhttp/cookiejar"
	"github.com/saucesteals/mimic"
	"golang.org/x/net/publicsuffix"
)

const homeURL = "https://www.apartments.com/"

func FetchHTML(targetURL string) (string, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return "", err
	}

	transport, err := mimic.NewTransport(mimic.TransportOptions{
		Version:  "147.0.0.0",
		Brand:    mimic.BrandChrome,
		Platform: mimic.PlatformMac,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			ForceAttemptHTTP2:   true,
			TLSHandshakeTimeout: 20 * time.Second,
		},
	})
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   45 * time.Second,
	}

	primeReq, err := http.NewRequest(http.MethodGet, homeURL, nil)
	if err != nil {
		return "", err
	}
	setDocumentHeaders(primeReq, homeURL)

	primeResp, err := client.Do(primeReq)
	if err != nil {
		return "", err
	}
	_, _ = io.Copy(io.Discard, primeResp.Body)
	primeResp.Body.Close()

	pageReq, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return "", err
	}
	setDocumentHeaders(pageReq, homeURL)

	resp, err := client.Do(pageReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func setDocumentHeaders(req *http.Request, referer string) {
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", referer)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header[http.HeaderOrderKey] = []string{
		"host",
		"cache-control",
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"upgrade-insecure-requests",
		"user-agent",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-user",
		"sec-fetch-dest",
		"referer",
		"accept-encoding",
		"accept-language",
		"priority",
	}
}
