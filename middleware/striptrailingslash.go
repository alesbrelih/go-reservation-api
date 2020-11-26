package middleware

import (
	"net/http"
	"regexp"
	"strings"
)

var pathRegex *regexp.Regexp

func init() {
	pathRegex = regexp.MustCompile("^/.*/$")
}

func StripImplementation(path string) string {
	if !pathRegex.MatchString(path) {
		return path
	}
	return strings.TrimRight(path, "/")
}

func StripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = StripImplementation(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
