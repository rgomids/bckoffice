-------------------------------------------------
-- promoters
-------------------------------------------------
CREATE TABLE promoters (
  id              CHAR(26) PRIMARY KEY,    -- ULID
  full_name       TEXT NOT NULL,
  email           TEXT UNIQUE,
  phone           TEXT,
  document_id     TEXT UNIQUE,             -- CPF/CNPJ
  bank_account    JSONB,                   -- dados bancários (ag, conta, pix)
  created_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at      TIMESTAMPTZ
);

-------------------------------------------------
-- commission_contracts
-------------------------------------------------
CREATE TABLE commission_contracts (
  id              CHAR(26) PRIMARY KEY,    -- ULID
  promoter_id     CHAR(26) REFERENCES promoters(id),
  percentage      NUMERIC(5,2) NOT NULL CHECK (percentage >= 0),
  starts_at       DATE NOT NULL,
  ends_at         DATE,                    -- NULL = vigente
  created_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at      TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at      TIMESTAMPTZ
);