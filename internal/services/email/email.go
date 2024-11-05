package email

import (
	"net/smtp"

	"SavingBooks/config"
)

type SmtpServer struct {
	Host           string
	Port           string
	senderEmail    string
	senderPassword string
}

func NewSmtpServer(config *config.Configuration) *SmtpServer {
	return &SmtpServer{
		Host:           config.EmailHost,
		Port:           config.EmailPort,
		senderEmail:    config.EmailSender,
		senderPassword: config.EmailAppPassword,
	}
}
func (s *SmtpServer) Address() string {
	return s.Host + ":" + s.Port
}

func (s *SmtpServer) SendEmail(receiverEmail, resetPasswordLink string) error {
	subject := "Subject: Password Reset Request\n"
	from := "From: " + s.senderEmail + "\n"
	to := "To: " + receiverEmail + "\n"


	body := "Hello,\n\n" +
		"We received a request to reset your password. Use the code below to reset it:\n\n" +
		"Reset password: " + resetPasswordLink + "\n\n" +
		"If you didn't request this, you can ignore this email.\n\n" +
		"This token will expired in 10 minutes,\n" +
		"Best regards,\n" +
		"The Support Team"

	message := []byte(from + to + subject + "\n" + body)

	auth := smtp.PlainAuth("", s.senderEmail, s.senderPassword, s.Host)
	err := smtp.SendMail(s.Address(), auth, s.senderEmail, []string{receiverEmail}, message)
	if err != nil {
		return err
	}
	return nil
}

