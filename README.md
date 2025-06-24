# bckOffice

Monolito modular que reúne **Go (Chi)** no backend, **Next.js 14 (App Router)** no frontend
e **PostgreSQL 16** como banco principal. Todo o ciclo de vida roda em **Docker Compose**,
tanto em desenvolvimento quanto em produção.

---

## Pré-requisitos

- Docker + Docker Compose
- Go 1.24 (para desenvolvimento backend)
- Node 20 (para desenvolvimento frontend)
 - Make (opcional) para atalhos — rode `make help` para listar comandos

## Subir em desenvolvimento

```bash
cp .env.example .env
docker-compose up -d
````

* **Frontend**: [http://localhost:3000](http://localhost:3000)
* **Backend**:  [http://localhost:8080/healthz](http://localhost:8080/healthz)
* **Postgres**: `postgres://rgps:rgps_pass@localhost:5432/rgps_backoffice`

## Observabilidade

```bash
docker-compose -f infra/observability.yml up -d
```

Acesse:
- Prometheus → <http://localhost:9090>
- Grafana    → <http://localhost:3001>  (admin / admin)

## Backup
docker-compose -f infra/backup.yml up -d

Console MinIO: <http://localhost:9001> (rgpsadmin / rgpssecret)

---

## Estrutura de pastas

```
backend/   # código Go (cmd, internal, pkg)
frontend/  # Next.js + Tailwind + HeadlessUI
infra/     # docker-compose, observability
migrations/ # scripts SQL (golang-migrate)
docs/      # documentação complementar
```

---

## Comandos úteis

| Ação               | Comando                                   |
| ------------------ | ------------------------------------------ |
| Testes backend     | `make test`                                |
| Lint backend       | `make lint`                                |
| Lint frontend      | `npm run lint`                             |
| Build imagens      | `make build`                               |
| Criar migration    | `make migrate-create name=exemplo`         |
| Aplicar migrations | `make migrate-up`                          |
| Desfazer migrations| `make migrate-down mnum=1`                 |
| Forçar versão      | `make migrate-up-force version=<N>`        |

---

## Roadmap resumido

* [ ] Autenticação + RBAC
* [ ] Módulo Clientes
* [ ] Serviços & Contratos
* [ ] Financeiro & Comissões
* [ ] Observabilidade Prometheus/Grafana

---

## Contribuindo

1. Rode `gofmt -w .` dentro do diretório `backend/` antes de commitar qualquer mudança em Go.
2. Execute `npm run lint` no `frontend/` para garantir a qualidade do código JavaScript e TypeScript.

## Imagens no GHCR

Cada merge na branch `main` gera imagens `:latest` e `:<commit_sha>` publicadas em **ghcr.io**.

Para atualizar em produção:

```bash
docker-compose -f infra/docker-compose.yml.prod pull backend frontend
docker-compose -f infra/docker-compose.yml.prod up -d
```

Opcional: crie o secret `GHCR_PAT` se preferir token permanente.

## Deploy automático

Crie os secrets `DEPLOY_HOST`, `DEPLOY_USER` e `DEPLOY_SSH_KEY` no GitHub.
O servidor de destino deve possuir Docker e docker-compose, além do diretório
`/opt/rgps-backoffice/` com permissão para o usuário configurado.

Ao finalizar o workflow **Build & Push Images** com sucesso, o job
**Deploy to Production** copia a pasta `infra/` via SSH e reinicia os
containers em produção.

Para testar manualmente:
1. Realize merge na branch `main`.
2. Aguarde o sucesso do workflow **Build & Push Images**.
3. Verifique se o job **Deploy to Production** renovou os containers no host.

---

## Licença

MIT
