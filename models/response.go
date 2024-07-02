package models

type LoginResponse struct {
	Token string `json:"token"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
