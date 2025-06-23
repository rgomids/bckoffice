# --- build stage ---
FROM node:20-alpine AS builder
WORKDIR /app

# cache de dependências
COPY package.json package-lock.json* ./
RUN npm install --frozen-lockfile

# cópia do restante do código
COPY . .

# build de produção
RUN npm run build

# --- run stage ---
FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=development
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./
RUN npm install --development --frozen-lockfile

EXPOSE 3000
CMD ["npm", "start"]
