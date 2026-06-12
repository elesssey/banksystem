package service

import (
	"banksystem/internal/model"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"strconv"
)

type EmailService interface {
	SendVerificationCode(to, code string) error
	SendVerification(to string, transaction *model.Transaction) error
	GenerateVerificationCode() (string, error)
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}

func (s *emailService) GenerateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *emailService) SendVerificationCode(to, code string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	auth := smtp.PlainAuth("", user, pass, host)

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: Подтверждение регистрации\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			"Ваш код подтверждения: " + code + "\r\n" +
			"Код действует 10 минут.\r\n",
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
}

func (s *emailService) SendVerification(to string, transaction *model.Transaction) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	auth := smtp.PlainAuth("", user, pass, host)

	yesURL := "http://localhost:8180/transaction/Confirm?id=" + strconv.Itoa(transaction.Id)
	noURL := "http://localhost:8180/transaction/Cancel?id=" + strconv.Itoa(transaction.Id)

	body := "Карта с номером " + transaction.SourseAccountNumber +
		" перевела " + fmt.Sprintf("%d", transaction.Amount) +
		" номеру " + transaction.DestinationAccountNumber + ".\r\n\r\n" +
		"Подтвердить операцию (да): " + yesURL + "\r\n" +
		"Отменить операцию (нет): " + noURL + "\r\n"

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: Уведомление о переводе\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body,
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
}
