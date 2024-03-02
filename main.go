package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{}

	r := chi.NewRouter()

	r.Mount("/api", api(apiCfg))
	r.Mount("/admin", admin(apiCfg))

	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port:%s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func api(cfg *apiConfig) chi.Router {
	r := chi.NewRouter()

	r.Get("/healthz", handlerReadiness)
	r.HandleFunc("/reset", cfg.handlerReset)

	return r
}

func admin(cfg *apiConfig) chi.Router {
	r := chi.NewRouter()

	r.Get("/metrics", cfg.handlerMetrics)

	return r
}
