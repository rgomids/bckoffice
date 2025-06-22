-------------------------------------------------
-- users
-------------------------------------------------
CREATE TABLE users (
  id           CHAR(26) PRIMARY KEY,        -- ULID
  email        TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  full_name    TEXT NOT NULL,
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at   TIMESTAMPTZ
);

-------------------------------------------------
-- roles
-------------------------------------------------
CREATE TABLE roles (
  id           CHAR(26) PRIMARY KEY,        -- ULID
  name         TEXT UNIQUE NOT NULL,
  description  TEXT,
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at   TIMESTAMPTZ
);

-------------------------------------------------
-- user_roles (N-to-N)
-------------------------------------------------
CREATE TABLE user_roles (
  user_id CHAR(26) REFERENCES users(id),
  role_id CHAR(26) REFERENCES roles(id),
  PRIMARY KEY (user_id, role_id)
);