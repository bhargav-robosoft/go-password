package entity

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type GenerateOtpRequest struct {
	Email string `json:"email" binding:"required"`
	Type  string `json:"type" binding:"required"`
}

type VerifyOtpRequest struct {
	Email string `json:"email" binding:"required"`
	Otp   string `json:"otp" binding:"required"`
}

type ResetPasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
