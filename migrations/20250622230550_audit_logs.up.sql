-------------------------------------------------
-- audit_logs
-------------------------------------------------
CREATE TABLE audit_logs (
  id             CHAR(26) PRIMARY KEY,        -- ULID
  user_id        CHAR(26) REFERENCES users(id),
  entity_name    TEXT NOT NULL,               -- ex.: 'customers'
  entity_id      CHAR(26),                    -- ULID da entidade
  action         TEXT NOT NULL CHECK (action IN (
                    'insert', 'update', 'delete'
                  )),
  diff           JSONB,                       -- fields changed (para update)
  ip_address     INET,
  user_agent     TEXT,
  geo_info       JSONB,                       -- {country, region, city, lat, lon}
  created_at     TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_name, entity_id);
CREATE INDEX idx_audit_logs_user   ON audit_logs (user_id);