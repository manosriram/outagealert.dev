package template

type Response struct {
	Message string
	Error   string
}

type RegisterSuccessResponse struct {
	Response
	Email string
}

type ForgotPasswordSuccessResponse struct {
	Response
	Email string
}
