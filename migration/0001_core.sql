-- +migrate Up

-- Tabela de usuários
CREATE TABLE usuarios (
  id           BIGSERIAL PRIMARY KEY,
  email        TEXT UNIQUE NOT NULL,
  senha_hash   TEXT NOT NULL,
  nome         TEXT NOT NULL,
  criado_em    TIMESTAMPTZ DEFAULT now() NOT NULL,
  atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
  apagado_em   TIMESTAMPTZ
);

-- Tabela de papéis (RBAC)
CREATE TABLE roles (
  id           BIGSERIAL PRIMARY KEY,
  nome         TEXT UNIQUE NOT NULL,
  descricao    TEXT,
  criado_em    TIMESTAMPTZ DEFAULT now() NOT NULL,
  atualizado_em TIMESTAMPTZ DEFAULT now() NOT NULL,
  apagado_em   TIMESTAMPTZ
);

-- Relação N-para-N usuário ↔ papel
CREATE TABLE user_roles (
  usuario_id BIGINT REFERENCES usuarios(id),
  role_id    BIGINT REFERENCES roles(id),
  PRIMARY KEY (usuario_id, role_id)
);

-- +migrate Down
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS usuarios;
