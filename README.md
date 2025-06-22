# RCM Tech — Backoffice

Monolito modular que reúne **Go (Chi)** no backend, **Next.js 14 (App Router)** no frontend
e **PostgreSQL 16** como banco principal. Todo o ciclo de vida roda em **Docker Compose**,
tanto em desenvolvimento quanto em produção.

---

## Pré-requisitos

- Docker + Docker Compose
- Go 1.22 (para desenvolvimento backend)
- Node 20 (para desenvolvimento frontend)
- Make (opcional) para atalhos

## Subir em desenvolvimento

```bash
cp .env.example .env
docker-compose -f infra/docker-compose.yml up -d --build
````

* **Frontend**: [http://localhost:3000](http://localhost:3000)
* **Backend**:  [http://localhost:8080/healthz](http://localhost:8080/healthz)
* **Postgres**: `postgres://rcm:rcm_pass@localhost:5432/rcm_backoffice`

---

## Estrutura de pastas

```
backend/   # código Go (cmd, internal, pkg)
frontend/  # Next.js + Tailwind + HeadlessUI
infra/     # docker-compose, observability
migration/ # scripts SQL (golang-migrate)
docs/      # documentação complementar
```

---

## Comandos úteis

| Ação              | Comando                 |
| ----------------- | ----------------------- |
| Testes backend    | `go test ./...`         |
| Lint frontend     | `npm run lint`          |
| Build backend     | `make build` (em breve) |
| Atualizar deps Go | `go get -u`             |

---

## Roadmap resumido

* [ ] Autenticação + RBAC
* [ ] Módulo Clientes
* [ ] Serviços & Contratos
* [ ] Financeiro & Comissões
* [ ] Observabilidade Prometheus/Grafana

---

## Licença

MIT
