package middlewares

import (
	"net/http"

	"saaa-api/internal/core/context"
)

const (
	// EN english language
	EN = "en"

	// TH thai language
	TH = "th"
)

// AcceptLanguage header Accept-Language
func AcceptLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		if lang != EN && lang != TH {
			lang = EN
		}
		context.SetLanguage(r, lang)
		next.ServeHTTP(w, r)
	})
}
