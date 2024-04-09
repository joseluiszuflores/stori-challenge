package email

import (
	"github.com/go-mail/mail"
	"github.com/golang/glog"
	"os"
)

type SMTPService struct {
	host         string
	port         int
	username     string
	password     string
	from         string
	templatePath string
}

const Subject = "Total Balance"

func (s *SMTPService) getTemplate() (string, error) {
	data, err := os.ReadFile(s.templatePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *SMTPService) Send(destination, name string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", destination)
	m.SetAddressHeader("Cc", destination, name)
	m.SetHeader("Subject", Subject)
	tmpl, err := s.getTemplate()
	if err != nil {
		return err
	}
	m.SetBody("text/html", tmpl)
	d := mail.NewDialer(s.host, s.port, s.username, s.password)
	if err := d.DialAndSend(m); err != nil {
		glog.Errorf("error sending the email [%s]", err)
		return err
	}

	return nil
}
