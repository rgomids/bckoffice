package customer

import "time"

// Address representa um endereco vinculado a um cliente.
type Address struct {
	ID          string     `db:"id" json:"id"`
	CustomerID  string     `db:"customer_id" json:"customerID"`
	AddressType string     `db:"address_type" json:"addressType"`
	Street      string     `db:"street" json:"street"`
	Number      string     `db:"number" json:"number"`
	Complement  string     `db:"complement" json:"complement"`
	District    string     `db:"district" json:"district"`
	City        string     `db:"city" json:"city"`
	State       string     `db:"state" json:"state"`
	PostalCode  string     `db:"postal_code" json:"postalCode"`
	Country     string     `db:"country" json:"country"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}
