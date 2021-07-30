package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/redhatinsights/internal/endpoints"
)


func lubdub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("lubdub"))
}


func main() {


	r := chi.NewRouter()
	mr := chi.NewRouter()
	sub := chi.NewRouter()

	// Mount the root of the api router on /api/v1
	r.Mount("/api/v1/", sub)
	r.Get("/", lubdub)

	// Mount the metrics handler on /metrics
	mr.Get("/", lubdub)
	mr.Handle("/metrics", promhttp.Handler())

	sub.Get("/", lubdub)
	sub.Get("/payloads", endpoints.Payloads)
	sub.Get("/payloads/{request_id}", endpoints.SinglePayload)
	sub.Get("/statuses", endpoints.Statuses)
	sub.Get("/health", endpoints.Health)

	srv := http.Server{
		Addr:	":8080",
		Handler: r,
	}

	msrv := http.Server{
		Addr:	":8081",
		Handler: mr,
	}

	go func() {

		if err := msrv.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}