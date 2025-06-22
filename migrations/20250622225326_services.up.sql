-------------------------------------------------
-- services
-------------------------------------------------
CREATE TABLE services (
  id           CHAR(26) PRIMARY KEY,          -- ULID
  name         TEXT UNIQUE NOT NULL,
  description  TEXT,
  base_price   NUMERIC(12,2) NOT NULL DEFAULT 0,
  is_active    BOOLEAN NOT NULL DEFAULT true,
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at   TIMESTAMPTZ
);