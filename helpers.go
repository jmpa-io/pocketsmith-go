package pocketsmith

import (
	"net/http"
	"regexp"
)

// the regex used to identify a rel in a http.Header.
var rgxRel = regexp.MustCompile(`<(.+?)>;\s*rel="(.+?)"`)

// getHeader check for the given key in the given headers and tries to
// extract the value attached to that header. It will either return the
// found value, or an empty string, if not found.
func getHeader(headers http.Header, key string) string {
	for _, link := range headers["Link"] {
		for _, m := range rgxRel.FindAllStringSubmatch(link, -1) {
			if len(m) != 3 {
				continue
			}
			if m[2] == key {
				return m[1]
			}
		}
	}
	return ""
}
