package auth

// CredentialsInput representa o payload de login.
type CredentialsInput struct {
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"admin123"`
}

// AuthResponse define a resposta contendo o token JWT.
type AuthResponse struct {
	Token string `json:"token"`
}
