-------------------------------------------------
-- contracts
-------------------------------------------------
CREATE TABLE contracts (
  id            CHAR(26) PRIMARY KEY,             -- ULID
  customer_id   CHAR(26) REFERENCES customers(id),
  service_id    CHAR(26) REFERENCES services(id),
  promoter_id   CHAR(26) REFERENCES promoters(id),
  value_total   NUMERIC(12,2) NOT NULL,
  start_date    DATE NOT NULL,
  end_date      DATE,
  status        TEXT NOT NULL DEFAULT 'active' CHECK (status IN (
                  'active', 'suspended', 'closed', 'cancelled'
                )),
  created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at    TIMESTAMPTZ
);

-------------------------------------------------
-- contract_attachments
-------------------------------------------------
CREATE TABLE contract_attachments (
  id            CHAR(26) PRIMARY KEY,             -- ULID
  contract_id   CHAR(26) REFERENCES contracts(id),
  file_name     TEXT NOT NULL,
  storage_url   TEXT NOT NULL,                    -- S3/minio path
  mime_type     TEXT,
  size_bytes    BIGINT,
  created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at    TIMESTAMPTZ
);