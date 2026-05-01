package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/fhttp/cookiejar"
	"github.com/saucesteals/mimic"
	"golang.org/x/net/publicsuffix"
)

const (
	homeURL   = "https://www.apartments.com/"
	targetURL = "https://www.apartments.com/off-campus-housing/ca/san-jose/san-jose-state-university/"
	outputDir = "ldjson"
)

func main() {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatalf("create cookie jar: %v", err)
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
		log.Fatalf("create mimic transport: %v", err)
	}

	client := &http.Client{
		Transport: transport,
		Jar:       jar,
		Timeout:   45 * time.Second,
	}

	primeReq, err := http.NewRequest(http.MethodGet, homeURL, nil)
	if err != nil {
		log.Fatalf("build homepage request: %v", err)
	}
	setDocumentHeaders(primeReq, homeURL)

	primeResp, err := client.Do(primeReq)
	if err != nil {
		log.Fatalf("prime apartments.com cookies: %v", err)
	}
	_, _ = io.Copy(io.Discard, primeResp.Body)
	primeResp.Body.Close()

	pageReq, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		log.Fatalf("build target request: %v", err)
	}
	setDocumentHeaders(pageReq, homeURL)

	resp, err := client.Do(pageReq)
	if err != nil {
		log.Fatalf("fetch apartments.com target: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read response: %v", err)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("create output dir: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Fatalf("parse html: %v", err)
	}

	count := 0
	doc.Find(`script[type="application/ld+json"]`).Each(func(i int, sel *goquery.Selection) {
		content := strings.TrimSpace(sel.Text())
		if content == "" {
			return
		}

		path := filepath.Join(outputDir, fmt.Sprintf("%d.json", count))
		if err := os.WriteFile(path, []byte(content+"\n"), 0o644); err != nil {
			log.Fatalf("write %s: %v", path, err)
		}
		fmt.Println(path)
		count++
	})

	fmt.Printf("wrote %d ld+json files\n", count)
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
