-------------------------------------------------
-- tags
-------------------------------------------------
CREATE TABLE tags (
  id           CHAR(26) PRIMARY KEY,    -- ULID
  name         TEXT UNIQUE NOT NULL,
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at   TIMESTAMPTZ
);

-------------------------------------------------
-- entity_tags  (associa qualquer tabela a uma tag)
-------------------------------------------------
CREATE TABLE entity_tags (
  id           CHAR(26) PRIMARY KEY,    -- ULID
  entity_name  TEXT NOT NULL,           -- ex.: 'customers', 'contracts'
  entity_id    CHAR(26) NOT NULL,       -- ULID da entidade
  tag_id       CHAR(26) REFERENCES tags(id),
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at   TIMESTAMPTZ
);

-------------------------------------------------
-- notes  (anotações internas com histórico)
-------------------------------------------------
CREATE TABLE notes (
  id           CHAR(26) PRIMARY KEY,    -- ULID
  entity_name  TEXT NOT NULL,
  entity_id    CHAR(26) NOT NULL,
  author_id    CHAR(26) REFERENCES users(id),
  text         TEXT NOT NULL,
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL
);