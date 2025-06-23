package customer

// AddressInput representa os dados de endereco recebidos na criacao.
type AddressInput struct {
	AddressType string `json:"address_type" validate:"required"`
	Street      string `json:"street" validate:"required"`
	Number      string `json:"number"`
	Complement  string `json:"complement"`
	District    string `json:"district"`
	City        string `json:"city" validate:"required"`
	State       string `json:"state" validate:"required"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
}

// CreateCustomerInput define o payload para criacao de clientes.
type CreateCustomerInput struct {
	LegalName  string         `json:"legal_name" validate:"required"`
	TradeName  string         `json:"trade_name"`
	DocumentID string         `json:"document_id" validate:"required"`
	Email      string         `json:"email" validate:"omitempty,email"`
	Phone      string         `json:"phone"`
	Addresses  []AddressInput `json:"addresses" validate:"required,min=1,dive"`
}
