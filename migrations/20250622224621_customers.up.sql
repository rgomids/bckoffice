-------------------------------------------------
-- customers
-------------------------------------------------
CREATE TABLE customers (
  id           CHAR(26) PRIMARY KEY,              -- ULID
  legal_name   TEXT NOT NULL,                     -- Razão Social ou Nome
  trade_name   TEXT,                              -- Nome fantasia, opcional
  document_id  TEXT UNIQUE,                       -- CPF/CNPJ
  email        TEXT,
  phone        TEXT,
  promoter_id  CHAR(26),                          -- FK futura (promoters)
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at   TIMESTAMPTZ
);

-------------------------------------------------
-- addresses
-------------------------------------------------
CREATE TABLE addresses (
  id           CHAR(26) PRIMARY KEY,              -- ULID
  customer_id  CHAR(26) REFERENCES customers(id),
  address_type TEXT NOT NULL DEFAULT 'billing',   -- billing | shipping | other
  street       TEXT NOT NULL,
  number       TEXT,
  complement   TEXT,
  district     TEXT,
  city         TEXT NOT NULL,
  state        TEXT NOT NULL,
  postal_code  TEXT,
  country      TEXT DEFAULT 'BR',
  created_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at   TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at   TIMESTAMPTZ
);