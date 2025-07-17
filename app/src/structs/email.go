package structs

type SendEmailEvent struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func (s *SendEmailEvent) SendEmail() {
}
