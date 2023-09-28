package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/christopherplain/chirpy/internal/api"
	"github.com/christopherplain/chirpy/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	const filePathRoot = "."
	const port = "8080"
	const dbPath = "database.json"

	godotenv.Load()

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		model.ResetDB(dbPath)
	}

	db, err := model.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := api.ApiConfig{
		DB:             db,
		FileserverHits: 0,
		JwtSecret:      os.Getenv("JWT_SECRET"),
		PolkaKey:       os.Getenv("POLKA_KEY"),
	}

	router := chi.NewRouter()
	fsHandler := apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Get("/reset", apiCfg.HandleReset)
	apiRouter.Get("/chirps", apiCfg.HandleGetChirps)
	apiRouter.Get("/chirps/{id}", apiCfg.HandleGetChirp)
	apiRouter.Delete("/chirps/{id}", apiCfg.HandleDeleteChirp)
	apiRouter.Post("/chirps", apiCfg.HandlePostChirp)
	apiRouter.Post("/login", apiCfg.HandleUserLogin)
	apiRouter.Post("/polka/webhooks", apiCfg.HandlePolkaWebhook)
	apiRouter.Post("/refresh", apiCfg.HandleRefresh)
	apiRouter.Post("/revoke", apiCfg.HandleRevoke)
	apiRouter.Post("/users", apiCfg.HandlePostUser)
	apiRouter.Put("/users", apiCfg.HandlePutUser)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.HandleMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
