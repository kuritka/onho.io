package rest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			logger.Info().
				Str("req_id", middleware.GetReqID(r.Context())).
				Str("uri", fmt.Sprintf("%s %s://%s%s %s from %s",
					r.Method, scheme, r.Host, r.RequestURI, r.Proto, r.RemoteAddr)).
				Str("user_agent", r.UserAgent()).
				Msg("parsing request")

			var h []string
			for k, v := range r.Header {
				h = append(h, fmt.Sprintf("%s:%s", k, v))
			}
			logger.Info().
				Str("header", string(strings.Join(h, "; "))).
				Msg("parsing header")

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}


// SecurityHeaderMiddleware injects security relevant headers
func SecurityHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// block MIME types sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// set CSP to default fetches to self origin
		// TODO: define proper CSP header
		//w.Header().Set("Content-Security-Policy", "default-src 'self'")
		//w.Header().Set("Content-Security-Policy", "style-src https://fonts.googleapis.com")
		// set XSS if CSP is not supported
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// set HSTS here, use SSL/TLS transport strictly
		w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		w.Header().Add("Strict-Transport-Security", "includeSubDomains")
		// response varies depending on `Origin` header
		w.Header().Set("Vary", "Origin")

		next.ServeHTTP(w, r)
	})
}



// CacheControlMiddleware extends header field Cache-Control with value `max-age=age`
func CacheControlMiddleware(age time.Duration) func(http.Handler) http.Handler {
	maxAge := fmt.Sprintf("max-age=%d", int64(age.Seconds()))
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			// define caching policy
			w.Header().Set("Cache-Control", maxAge)
			w.Header().Add("Cache-Control", "must-revalidate")

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}

// SetHTTPHeader middleware sets a header for the certain requested URL paths
// if urls slice is nil it just sets a header
func SetHTTPHeader(headName, headVal string, urls []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			if urls == nil {
				w.Header().Set(headName, headVal)
			} else {
				for _, v := range urls {
					if v == r.URL.Path {
						w.Header().Set(headName, headVal)
					}
				}
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hfn)
	}
}
