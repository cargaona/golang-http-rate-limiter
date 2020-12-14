package main

import (
	"log"
	"net/http"
)

var limiter = NewIPRateLimiter(490, 200)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ping)

	if err := http.ListenAndServe(":8889", limitMiddleware(mux)); err != nil {
		log.Fatal("Unable to start server: %s", err.Error())
	}
}

func ping(responseWriter http.ResponseWriter, r *http.Request) {
	responseWriter.Write([]byte("Finished"))
//	fmt.Fprintf(responseWriter, "Finished!")
}

func limitMiddleware(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		limiter := limiter.GetLimiter(request.RemoteAddr)
		if !limiter.Allow() {
			http.Error(responseWriter, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		mux.ServeHTTP(responseWriter, request)
	})
}
