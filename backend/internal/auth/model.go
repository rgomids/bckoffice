package auth

// CredentialsInput representa o payload de login.
type CredentialsInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse define a resposta contendo o token JWT.
type AuthResponse struct {
	Token string `json:"token"`
}
