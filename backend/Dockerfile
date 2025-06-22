# --- build stage ---
FROM golang:1.22-alpine AS builder
WORKDIR /app

# cache das dependências
COPY go.mod go.sum ./
RUN go mod download

# cópia do restante do código
COPY . .

# build estático
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# --- run stage ---
FROM gcr.io/distroless/static
COPY --from=builder /server /server
ENV PORT=8080
EXPOSE 8080
CMD ["/server"]
