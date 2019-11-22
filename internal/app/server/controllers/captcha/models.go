package captcha

type ReqEmail struct {
	Email string `json:"email" validate:"required,email" binding:"required,email"`
	For   string `json:"for" validate:"required" binding:"required"`
}
