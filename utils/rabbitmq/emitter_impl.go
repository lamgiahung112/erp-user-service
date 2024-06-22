package rabbitmq

type LoginOtpEmailPayload struct {
	ToAddress string `json:"to_address"`
	Title     string `json:"title"`
	Otp       string `json:"otp"`
	Username  string `json:"username"`
}

type VerifyAccountEmailPayload struct {
	ToAddress   string `json:"to_address"`
	Title       string `json:"title"`
	RedirectUrl string `json:"redirect_url"`
	Username    string `json:"username"`
}

func (e *EventEmitter) SendLoginOtpEmail(payload *LoginOtpEmailPayload) error {
	return e.pushEmailRequest(LoginOTP, payload)
}

func (e *EventEmitter) SendVerifyAccountEmail(payload *VerifyAccountEmailPayload) error {
	return e.pushEmailRequest(VerifyAccount, payload)
}
