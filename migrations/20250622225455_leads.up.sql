-------------------------------------------------
-- leads
-------------------------------------------------
CREATE TABLE leads (
  id            CHAR(26) PRIMARY KEY,             -- ULID
  customer_id   CHAR(26) REFERENCES customers(id),
  promoter_id   CHAR(26) REFERENCES promoters(id),
  service_id    CHAR(26) REFERENCES services(id),
  status        TEXT NOT NULL CHECK (status IN (
                  'lead', 'qualified', 'proposal', 'contract'
                )),
  notes         TEXT,                              -- observações rápidas
  created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at    TIMESTAMPTZ
);