package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

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

	r := chi.NewRouter()

	// módulo Customers
	customer.RegisterRoutes(r, customerRepo)
	// módulo Services
	service.RegisterRoutes(r, serviceRepo)
	// módulo Promoters
	promoter.RegisterRoutes(r, promoterRepo)

	// rota simples de health-check
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("▶️  backend rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
