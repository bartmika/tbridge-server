package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/bartmika/tbridge-server/internal/utils"
)

// Middleware will split the full URL path into slash-sperated parts and save to
// the context to flow downstream in the app for this particular request.
func (h *ServerOperationImpl) URLProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Split path into slash-separated parts, for example, path "/foo/bar"
		// gives p==["foo", "bar"] and path "/" gives p==[""]. Our API starts with
		// "/api", as a result we will start the array slice at "2".
		p := strings.Split(r.URL.Path, "/")[2:]

		// log.Println(p) // For debugging purposes only.

		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "url_split", p)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

func (h *ServerOperationImpl) BearerProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract our auth header array.
		reqToken := r.Header.Get("Authorization")

		// // For debugging purposes.
		// log.Println("BearerProcessorMiddleware | reqToken:", reqToken)

		// Before running our Bearer middleware we need to confirm there is an
		// an `Authorization` header to run our middleware. This is an important
		// step!
		if reqToken != "" && strings.Contains(reqToken, "undefined") == false {

			// Special thanks to "poise" via https://stackoverflow.com/a/44700761
			splitToken := strings.Split(reqToken, "Bearer ")
			if len(splitToken) < 2 {
				http.Error(w, "not properly formatted authorization header", http.StatusBadRequest)
				return
			}

			reqToken = splitToken[1]

			sessionUuid, err := utils.ProcessBearerToken(h.hmacSecret, reqToken)
			// log.Println(sessionUuid) // For debugging purposes only.

			if err == nil {
				// Update our context to save our Bearer token content information.
				ctx = context.WithValue(ctx, "is_authorized", true)
				ctx = context.WithValue(ctx, "session_uuid", sessionUuid)

				// Flow to the next middleware with our Bearer token saved.
				fn(w, r.WithContext(ctx))
				return
			}

			// The following code will lookup the URL path in a whitelist and
			// if the visited path matches then we will skip any token errors.
			// We do this because a majority of API endpoints are protected
			// by authorization.

			urlSplit := ctx.Value("url_split").([]string)
			skipPath := map[string]bool{
				"register":      true,
				"login":         true,
				"refresh-token": true,
			}

			// DEVELOPERS NOTE:
			// If the URL cannot be split into the size we want then skip running
			// this middleware.
			if len(urlSplit) >= 2 {
				if skipPath[urlSplit[1]] {
					log.Println("BearerProcessorMiddleware | ProcessBearer | Skipping expired or error token")
				} else {
					log.Println("BearerProcessorMiddleware | ProcessBearer | err", err, "for reqToken:", reqToken)
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
			}
		}

		// Flow to the next middleware without anything done.
		ctx = context.WithValue(ctx, "is_authorized", false)
		fn(w, r.WithContext(ctx))
	}
}

func (h *ServerOperationImpl) AuthorizationMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get our authorization information.
		isAuthorized, ok := ctx.Value("is_authorized").(bool)
		if ok && isAuthorized {
			sessionUuid := ctx.Value("session_uuid").(string)

            // DEVELOPERS NOTE:
			// This is for added security.
			if sessionUuid != h.sessionUuid {
				http.Error(w, "Session expired or does not exist - please log in again", http.StatusUnauthorized)
				return
			}
		}

		fn(w, r.WithContext(ctx))
	}
}

func (h *ServerOperationImpl) IPAddressMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the IPAddress. Code taken from: https://stackoverflow.com/a/55738279
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		// Save our IP address to the context.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "IPAddress", IPAddress)
		fn(w, r.WithContext(ctx)) // Flow to the next middleware.
	}
}

// The purpose of this middleware is to return a `401 unauthorized` error if
// the user is not authorized and visiting a protected URL.
func (h *ServerOperationImpl) ProtectedURLsMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// The following code will lookup the URL path in a whitelist and
		// if the visited path matches then we will skip URL protection.
		// We do this because a majority of API endpoints are protected
		// by authorization.

		urlSplit := ctx.Value("url_split").([]string)
		skipPath := map[string]bool{
			"version":       true,
			"register":      true,
			"login":         true,
			"refresh-token": true,
		}

		// DEVELOPERS NOTE:
		// If the URL cannot be split into the size we want then skip running
		// this middleware.
		if len(urlSplit) <= 1 {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		if skipPath[urlSplit[1]] {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
		} else {
			// Get our authorization information.
			isAuthorized, ok := ctx.Value("is_authorized").(bool)

			// Either accept continuing execution or return 401 error.
			if ok && isAuthorized {
				fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			} else {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
	}
}

func (h *ServerOperationImpl) ChainMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	// Attach our middleware handlers here. Please note that all our middleware
	// will start from the bottom and proceed upwards.
	// Ex: `URLProcessorMiddleware` will be executed first and
	//     `AuthorizationMiddleware` will be executed last.
	fn = h.ProtectedURLsMiddleware(fn)
	fn = h.IPAddressMiddleware(fn)
	fn = h.AuthorizationMiddleware(fn) // Note: Must be above `BearerProcessorMiddleware`.
	fn = h.BearerProcessorMiddleware(fn)
	fn = h.URLProcessorMiddleware(fn)

	return func(w http.ResponseWriter, r *http.Request) {
		// Flow to the next middleware.
		fn(w, r)
	}
}
