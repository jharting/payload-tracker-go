package main

import (
	// "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// "github.com/redhatinsights/payload-tracker-go/models"
	"github.com/redhatinsights/payload-tracker-go/pkg/db"
	"github.com/redhatinsights/payload-tracker-go/internal/config"
	"github.com/redhatinsights/payload-tracker-go/internal/endpoints"
)

func lubdub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("lubdub"))
}

func main() {

	cfg := config.Get()

	db.DbConnect()

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
		Addr:    ":"+cfg.PublicPort,
		Handler: r,
	}

	msrv := http.Server{
		Addr:    ":"+cfg.MetricsPort,
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


	// db test
	// data := db.DB.Take(&models.Payloads)
	// fmt.Printf("%v", data)
}
