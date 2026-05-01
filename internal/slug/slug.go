package slug

import (
	"net/url"
	"path"
	"strings"
)

func FromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "apartments"
	}

	base := path.Base(strings.TrimSuffix(u.Path, "/"))
	if base == "." || base == "/" || base == "" {
		return "apartments"
	}
	return base
}
