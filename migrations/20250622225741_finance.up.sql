-------------------------------------------------
-- accounts_receivable
-------------------------------------------------
CREATE TABLE accounts_receivable (
  id             CHAR(26) PRIMARY KEY,            -- ULID
  contract_id    CHAR(26) REFERENCES contracts(id),
  due_date       DATE NOT NULL,
  amount         NUMERIC(12,2) NOT NULL,
  status         TEXT NOT NULL DEFAULT 'open' CHECK (status IN (
                    'open', 'paid', 'overdue', 'cancelled'
                  )),
  paid_at        TIMESTAMPTZ,
  created_at     TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at     TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at     TIMESTAMPTZ
);

-------------------------------------------------
-- commissions
-------------------------------------------------
CREATE TABLE commissions (
  id             CHAR(26) PRIMARY KEY,            -- ULID
  contract_id    CHAR(26) REFERENCES contracts(id),
  promoter_id    CHAR(26) REFERENCES promoters(id),
  amount         NUMERIC(12,2) NOT NULL,
  approved       BOOLEAN NOT NULL DEFAULT false,
  approved_by    CHAR(26) REFERENCES users(id),
  approved_at    TIMESTAMPTZ,
  created_at     TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at     TIMESTAMPTZ DEFAULT now() NOT NULL,
  deleted_at     TIMESTAMPTZ
);