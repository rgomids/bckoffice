package audit

import (
	"encoding/json"
	"time"
)

// AuditLog representa um registro de alteração na aplicação.
// Campos diff e geo_info são armazenados como JSONB.
type AuditLog struct {
	ID         string          `db:"id" json:"id"`
	UserID     string          `db:"user_id" json:"userId"`
	EntityName string          `db:"entity_name" json:"entityName"`
	EntityID   string          `db:"entity_id" json:"entityId"`
	Action     string          `db:"action" json:"action"`
	Diff       json.RawMessage `db:"diff" json:"diff,omitempty"`
	IPAddress  string          `db:"ip_address" json:"ipAddress"`
	UserAgent  string          `db:"user_agent" json:"userAgent"`
	GeoInfo    json.RawMessage `db:"geo_info" json:"geoInfo"`
	CreatedAt  time.Time       `db:"created_at" json:"createdAt"`
}
