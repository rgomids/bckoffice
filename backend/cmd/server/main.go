package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/smithl4b/rcm.backoffice/internal/auth"
	"github.com/smithl4b/rcm.backoffice/internal/contract"
	"github.com/smithl4b/rcm.backoffice/internal/customer"
	"github.com/smithl4b/rcm.backoffice/internal/promoter"
	"github.com/smithl4b/rcm.backoffice/internal/service"
)

func main() {

	dsn := os.Getenv("DB_DSN")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	customerRepo := customer.NewPostgresRepository(db)
	serviceRepo := service.NewPostgresRepository(db)
	promoterRepo := promoter.NewPostgresRepository(db)
	contractRepo := contract.NewPostgresRepository(db)
	authRepo := auth.NewPostgresRepository(db)

	r := chi.NewRouter()

	// rota publica de login
	auth.RegisterRoutes(r, authRepo)

	// rotas protegidas
	r.Group(func(pr chi.Router) {
		pr.Use(auth.AuthMiddleware)

		pr.Group(func(r chi.Router) {
			r.Use(auth.RequireRole("admin", "finance"))
			customer.RegisterRoutes(r, customerRepo)
		})

		service.RegisterRoutes(pr, serviceRepo)
		promoter.RegisterRoutes(pr, promoterRepo)
		contract.RegisterRoutes(pr, contractRepo)
	})

	// rota simples de health-check
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("▶️  backend rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
