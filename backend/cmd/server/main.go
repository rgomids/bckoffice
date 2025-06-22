package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()

    // rota simples de health-check
    r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
        w.Write([]byte("ok"))
    })

    log.Println("▶️  backend rodando em http://localhost:8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}
