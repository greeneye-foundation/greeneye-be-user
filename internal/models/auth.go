package models

type PasswordRecoveryRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required,e164"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}