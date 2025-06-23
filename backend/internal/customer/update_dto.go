package customer

// UpdateCustomerInput define o payload para atualizacao de clientes.
type UpdateCustomerInput struct {
	LegalName  string         `json:"legal_name" validate:"required"`
	TradeName  string         `json:"trade_name"`
	DocumentID string         `json:"document_id" validate:"required"`
	Email      string         `json:"email" validate:"omitempty,email"`
	Phone      string         `json:"phone"`
	Addresses  []AddressInput `json:"addresses" validate:"omitempty,dive"`
}
