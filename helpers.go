package pocketsmith

import (
	"net/http"
	"reflect"
	"regexp"
)

// Ahe regex used to identify a rel in a http.Header.
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

// A helper function to check if an interface has a value or not.
// https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1.
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array,
		reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
