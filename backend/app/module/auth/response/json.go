package response

type LoginResponse struct {
	Token     string `json:"token"`
	Type      string `json:"type"`
	ExpiresAt int64  `json:"expires_at"`
}

type RegisterResponse struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
}

type UserResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
