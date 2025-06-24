// @title        RCM Backoffice API
// @version      0.1
// @description  Endpoints do backoffice RCM Tech.
// @host         localhost:8080
// @BasePath     /
// @schemes      http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/rgomids/bckoffice/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/rgomids/bckoffice/internal/audit"
	"github.com/rgomids/bckoffice/internal/auth"
	"github.com/rgomids/bckoffice/internal/contract"
	"github.com/rgomids/bckoffice/internal/customer"
	"github.com/rgomids/bckoffice/internal/finance"
	"github.com/rgomids/bckoffice/internal/lead"
	"github.com/rgomids/bckoffice/internal/promoter"
	"github.com/rgomids/bckoffice/internal/service"
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
	leadRepo := lead.NewPostgresRepository(db)
	contractRepo := contract.NewPostgresRepository(db)
	financeRepo := finance.NewPostgresRepository(db)
	authRepo := auth.NewPostgresRepository(db)
	auditRepo := audit.NewPostgresRepository(db)
	geoSvc := audit.NewHttpGeoService()

	r := chi.NewRouter()
	corsMw := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Location"},
		AllowCredentials: true,
	})
	r.Use(corsMw.Handler)
	auditMw := audit.NewAuditMiddleware(auditRepo, geoSvc)
	r.Use(auditMw)
	r.Get("/docs/*", httpSwagger.WrapHandler)

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
		lead.RegisterRoutes(pr, leadRepo)
		contract.RegisterRoutes(pr, contractRepo)
		finance.RegisterRoutes(pr, financeRepo)
	})

	// rota simples de health-check
	// @Summary      Health check
	// @Tags         status
	// @Success      200  {string}  string  "ok"
	// @Router       /healthz [get]
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("▶️  backend rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
