# AGENTS.md

> **Objetivo**  
> Documentar os “agentes” (scripts ou jobs automatizados) que cuidam das rotinas
> de build, teste, deploy e manutenção do projeto **rcm-backoffice**.  
> Cada agente é independente, podendo ser invocado manualmente ou via CI.

| Agente | Diretório / entrypoint | Responsabilidade principal |
|--------|-----------------------|----------------------------|
| **CodexAssistant** | – | Recebe prompts curtos e gera código ou shell-scripts conforme as convenções deste repositório. |
| **BackendBuilder** | `backend/` – `make build-be` | Compilar o binário Go (multi-stage), rodar `go vet` e testes de unidade. |
| **FrontendBuilder** | `frontend/` – `make build-fe` | Gerar artefatos Next.js para produção e verificar ESLint/TypeScript. |
| **DBMigrator** | `migrations/` – `make migrate-*` | Criar e aplicar migrações SQL no PostgreSQL (create, up, down, force). |
| **Deployer** | `infra/docker-compose.yml` | Fazer `docker-compose pull && up -d` em staging/produção. |
| **ObservabilityWatcher** | `infra/observability.yml` | Subir stack Prometheus/Grafana e validar dashboards. |
| **BackupAgent** | `infra/backup/pg_dump.sh` | Executar backup diário do banco e enviar para bucket S3-compatível. |
| **CI** | `.github/workflows/ci.yml` | Rodar `go vet ./...` e `go test ./...` e build do frontend em cada push. |
| **ImagePublisher** | `.github/workflows/docker-build.yml` | Build e push das imagens backend/frontend para o GHCR a cada merge na `main`. |

### Boas práticas para novos agentes

1. **Isolamento** – cada agente deve ter script ou make-target próprio.  
2. **Logs claros** – prefixar saídas com o nome do agente.  
3. **Idempotência** – rodar o agente duas vezes não deve quebrar o estado.  
4. **Falha rápida** – exit code ≠ 0 e stack-trace legível.
5. **Documentação** – atualizar esta tabela ao adicionar ou alterar agente, atualizar README.md e AGENTS.md sempre que necessário
6. **Revisão de código** – ao propor modificações em arquivos, analisar o trecho alterado para identificar possíveis melhorias ou riscos antes de aplicar a mudança.
7. **Testes de unidade** – sempre criar e rodar testes unitários para qualquer alteração de código.
8. **Formatação** – utilizar ferramentas de formatação (ex.: `gofmt`, `prettier`) para garantir a qualidade e consistência do código.
9. **Gofmt obrigatório** – antes de realizar commits, execute `gofmt -w .` dentro do diretório `backend/` para manter o código Go devidamente formatado.

