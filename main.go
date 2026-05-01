package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	http "github.com/saucesteals/fhttp"
	"github.com/saucesteals/mimic"
)

const targetURL = "https://www.apartments.com/off-campus-housing/ca/san-jose/san-jose-state-university/"

const cookieHeader = `cb=1; cul=en-US; afe=%7B%22e%22%3Afalse%7D; fso=%7B%22e%22%3Afalse%7D; ss=1; ssot=1; _ga=GA1.2.623445881.1777578810; _gid=GA1.2.503430919.1777578810; gip=%7B%22Display%22%3A%22Chicago%2C%20IL%22%2C%22GeographyType%22%3A2%2C%22Address%22%3A%7B%22City%22%3A%22Chicago%22%2C%22CountryCode%22%3A%22US%22%2C%22State%22%3A%22IL%22%7D%2C%22Location%22%3A%7B%22Latitude%22%3A41.8337%2C%22Longitude%22%3A-87.7319%7D%2C%22IsPmcSearchByCityState%22%3Afalse%2C%22IsAreaTooFar%22%3Afalse%7D; s=; AKA_A2=A; bm_ss=ab8e18ef4e; sr=%7B%22Width%22%3A1303%2C%22Height%22%3A1130%2C%22PixelRatio%22%3A2%7D; akaalb_www_apartments_com_main=1777659640~op=ap_alb_aptsweb_prd_activity_logging:www_apartments_com_LAX|ap_alb_aptsweb_prd_reston_only:www_apartments_com_RESTON_CORE|ap_alb_aptsweb_prd_default:www_apartments_com_RESTON_CORE|~rv=33~m=www_apartments_com_LAX:0|www_apartments_com_RESTON_CORE:0|~os=0847b47fe1c72dfaedb786f1e8b4b630~id=3f3bb0d1c8f8713f9200666752f25774; bm_mi=00E45534008E8EC169C0EC6D1A7A6B4A~YAAQRckLFy+V4NidAQAAP/GO5B9h5Sp2RxECoJBlqprxL22xHgM7oW8WebPwPks0KEHUrdtl+njPXNOqFz/ZT5SL2a3058r1asmEbzBBIADmEt24G6rp8sh1WkgizTHtGcBr7gUb/RNDdlLW08Ohu2HctShVc9iHf7yRmPQKVE3Krh4jbLsjv2BjKcvzcNZdqu3QqfV8RHj7oah6i6Km9RaeIiyNBRwjUjuRx/uFdVDLS8l3mF/93Xi5Qr5VxBFosqCzTpOQEg5nLb1g7UFXn1hywdfJLSWAyUJqWGYGTb5XxcOgIZW+FotgABIUzwjf+K22cOQA/Fd4+wOowMjk4e7Izb8OZsm8Qq7XyC7MrOCHGpFSEZ+7QTBaAsX+6kKwojsqf1InoNv881c=~1; bm_sv=3BCCA0F27394DCD46913DC5C582D5721~YAAQRckLF6OX4NidAQAA7vaO5B+tFMeQ/7OhrL54YOpzoIDw3oCBiiFE5UVcVeHS/wBlwND2wAge1uMlY+2Bv0JMYqgCAsPJIKzbElUcP55sMOmTfxHFy1THmyz0zp95kxTGs0PlHzAjFS8MFNK+Gp8hokEKKLMR8/SEIICtGEJWhSzXu+2T1XV44YrV0HEcCf8Ib4o91PebV7pAI5A35eR7dU86uGnKFYNArgB7lWYTqFLAZzqi0BA9+hnY/v7wQBYHYAQ=~1; ak_bmsc=77EC6E0931EA54D4C406D3A7EB4D9970~000000000000000000000000000000~YAAQRckLF4C04didAQAAb42R5B8gqxgj1zhWTbqRNor3Y8fY7ssBSxiiLVQsGOXc9epwpM4ZPJAtOoAajuir6iYc/6CV+cVc3DZbl2JpLTCw4QKLqtcm9Y7+SDfJOuel7Z4c1tcPrx9l1kWsUVo7OK8c6MHCGk6rkadEiowwc76NIObD63DXH+yLaY2J3gkcNitORNt43bvULDU5u62RpG39yqWVKtxagZWjGXS2NTkMGLmB2uWZYn9ahwFBt4KPFz4MxlO+4cNtBNmRKKbpjTjJQFla4Mlin2UpqyFTNDpX2bVut61sSi721m6ssBd0lJEalONu8xb7vhiKs3gsQNO5S1ZlkjAUF9UI/JibDEbgyvjXyF54GX39D8qT/BqPPghrybBrw3Av4Vh53A6h/0oiKgWqoBlEyoxFtPHiN9zdCHKQVVfTykUfYlsftqDU46UIQqm1xURELsZfrodSY7J3kfrq1oGucPWRgxuZ6k6xEHJokg==; uat=%7B%22VisitorId%22%3A%2253d10dcc-b85c-435b-be70-61732d57a30c%22%2C%22VisitId%22%3A%2233935f09-aa5e-409b-b049-9c0127e80eb9%22%2C%22LastSearchId%22%3A%22567FAAEA-B9D9-47B2-9EE5-5A357506101D%22%7D; lsc=%7B%22Map%22%3A%7B%22BoundingBox%22%3A%7B%22LowerRight%22%3A%7B%22Latitude%22%3A37.31734%2C%22Longitude%22%3A-121.91217%7D%2C%22UpperLeft%22%3A%7B%22Latitude%22%3A37.35143%2C%22Longitude%22%3A-121.94465%7D%7D%2C%22CountryCode%22%3A%22US%22%7D%2C%22Geography%22%3A%7B%22ID%22%3A%22y9kdwss%22%2C%22Display%22%3A%22San%20Jose%20State%20University%20-%20San%20Jose%2C%20CA%20(University)%22%2C%22GeographyType%22%3A13%2C%22Address%22%3A%7B%22City%22%3A%22San%20Jose%22%2C%22CountryCode%22%3A%22USA%22%2C%22County%22%3A%22Santa%20Clara%22%2C%22PostalCode%22%3A%2295192%22%2C%22State%22%3A%22CA%22%2C%22StreetName%22%3A%22%22%2C%22Title%22%3A%22San%20Jose%20State%20University%22%2C%22Abbreviation%22%3A%22SJSU%22%2C%22MarketName%22%3A%22San%20Jose%22%2C%22DMA%22%3A%22San%20Francisco-Oakland-San%20Jose%2C%20CA%22%7D%2C%22Location%22%3A%7B%22Latitude%22%3A37.3358%2C%22Longitude%22%3A-121.881%7D%2C%22BoundingBox%22%3A%7B%22LowerRight%22%3A%7B%22Latitude%22%3A37.31495%2C%22Longitude%22%3A-121.84479%7D%2C%22UpperLeft%22%3A%7B%22Latitude%22%3A37.36148%2C%22Longitude%22%3A-121.90842%7D%7D%2C%22Radius%22%3A3%2C%22v%22%3A226%2C%22IsPmcSearchByCityState%22%3Afalse%2C%22IsAreaTooFar%22%3Afalse%7D%2C%22Listing%22%3A%7B%7D%2C%22Paging%22%3A%7B%7D%2C%22ResultSeed%22%3A481567%2C%22Options%22%3A0%2C%22CountryAbbreviation%22%3A%22US%22%2C%22ExcludeOptions%22%3A0%2C%22IsCoreGeoSearch%22%3Atrue%2C%22EnableAdvancedSearchSchoolRating%22%3A0%7D; csgp-origin=Reston-CSGP-CORE; bm_s=YAAQRckLF2Fo5didAQAAHtGZ5AXn86HkaNRGYIvpSj5Im/dOgJUgOnL1hJRoobBR/vVHC8nol6PGbom0bq3FYV2YJwEN7g8j4rmawpnIrkn+wCfFrX0YxH/1T9Sx0MbY4JrFCX7vm8Fg/berznCIRyNJaXgsJ2rwXrhfuF8P7jIoPDGXohEsKgHKCDfgSuoB56PJ7dJAsDbLPSs/b1lTUzkwFn37T09ncpvfc1kWnPR2QyFsFlT/+OBya/ztNVvgf01Lq3HTyzVbqzAth5eSz5jLyODA8hDc1hvpyV3DYkVIsxdjhwsJ5jACLyK9WBPJHAYceTyr7D0DTz6PXfF9bM5YkD0p5aVxN/x8f1YFTlwVcLxWNRaqgUQbSWCiakCldJv+TpSmcLs5Ng7fJDWAbr8MHaT7HlyNLI7/yzEdZQckyYXHEAMRyHFLdaJ9XC19GVdaBFpyc3X9lbF7TLi5QK8m0X12yBum/yHmXwgtgvQqUoi/axEHU119AxPePpMbaCdIm+rzQSv6qP2Dnn8v/ecQ8MLZPrWL3tJbqpKXgRyjhA4Vdx6sYUrB0J02yxZwLuc2yS40IHzV+9SNK9+dFXxm/OFZuMYZy3a/67zklAS9h98DLugbAdGdW7zyx/6Xazyryy6wg1WUxBgX/FGtQP7IubkJSvFqXC/lXJYb9oxH7lkp2rN/zT+gmjWeR48ip/oKjbC6+5kPn0XNkfWU+IAfpeZ0sgd8dwiYJrK1I2njOT8XOl+agVNRuLMlsuo/NlTGp6mv1e1RmwbL78JO0yd10G0jht8FXJMF21uc0ljX4+Zykz4QhQwkTnJu9GZ644l1IP46yntAjfhGL2gP0jXZYnR5O0ddpNMc+xb9S5DoIEzGOS+L5ELdEVvWcnf1Vt6yJUeDJq6Y; bm_so=4140DE72BFDC60C50B325DB126A8C109C7E44233D169789B64BBB75D52373D42~YAAQRckLF2Jo5didAQAAHtGZ5AeSkcCJ8rkbeitHuKlgkaryuDKvAR5XcYzqZRMtOolDUnUgetBmyzzYQcWrhNwW2aUCje/a4/0hxrVvMlioNR86rqxj/8R+RIhvIknaULQj31OoR2Lt9Do8cAKyuE+RoG9RidM2e8A4fSl1t8USirJKk+tCd4u3gxsU068Pb0okZY/hT+Mc76SYVmvNU2Hl4oKwO1dQJd1PrGC+FaCeAqck50ZfjIIm9PWKYqHn1uxSnWaVljJmlaWpF0PQRzhnFml3dPgjkckTmFFUweN7eBAiaLwSs2AjSza5UZQdNpoL3eOd8ZMWja/3mYQHFn4834h3O5tWf60HE2sZD4NmMmBoCdEe0ili8Kif3ePPrA1/rlg9QEMxyoyvGoSZr/AawE7WeKlXc4N/giNLhJioZnDiUWJevvf0bsJWvEWHrjg/Lck6CNeWMnEw0Ft0kWx+rjmdVbg=; bm_lso=4140DE72BFDC60C50B325DB126A8C109C7E44233D169789B64BBB75D52373D42~YAAQRckLF2Jo5didAQAAHtGZ5AeSkcCJ8rkbeitHuKlgkaryuDKvAR5XcYzqZRMtOolDUnUgetBmyzzYQcWrhNwW2aUCje/a4/0hxrVvMlioNR86rqxj/8R+RIhvIknaULQj31OoR2Lt9Do8cAKyuE+RoG9RidM2e8A4fSl1t8USirJKk+tCd4u3gxsU068Pb0okZY/hT+Mc76SYVmvNU2Hl4oKwO1dQJd1PrGC+FaCeAqck50ZfjIIm9PWKYqHn1uxSnWaVljJmlaWpF0PQRzhnFml3dPgjkckTmFFUweN7eBAiaLwSs2AjSza5UZQdNpoL3eOd8ZMWja/3mYQHFn4834h3O5tWf60HE2sZD4NmMmBoCdEe0ili8Kif3ePPrA1/rlg9QEMxyoyvGoSZr/AawE7WeKlXc4N/giNLhJioZnDiUWJevvf0bsJWvEWHrjg/Lck6CNeWMnEw0Ft0kWx+rjmdVbg=~1777656780449`

func main() {
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
		Timeout:   45 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		log.Fatalf("build request: %v", err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("cookie", cookieHeader)
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://www.apartments.com/")
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
		"cookie",
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("fetch apartments.com: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read response: %v", err)
	}

	fmt.Printf("status=%s\n", resp.Status)
	fmt.Printf("content-type=%s\n", resp.Header.Get("content-type"))
	fmt.Printf("content-encoding=%s\n", resp.Header.Get("content-encoding"))
	fmt.Printf("body-bytes=%d\n\n", len(body))

	preview := string(body)
	if len(preview) > 1500 {
		preview = preview[:1500]
	}
	fmt.Println(strings.TrimSpace(preview))
}
