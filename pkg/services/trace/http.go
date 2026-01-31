package trace

import (
	"net/http"
)

func (t *Tracer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := t.Context(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
