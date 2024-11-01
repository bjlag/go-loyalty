package middleware

import (
	"net/http"
	"strings"
)

func AllowMethods(methods ...string) func(next http.Handler) http.Handler {
	m := make(map[string]struct{}, len(methods))
	for _, method := range methods {
		m[strings.ToUpper(method)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := m[strings.ToUpper(r.Method)]; !ok {
				w.Header().Set("Allow", strings.Join(methods, ", "))
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
