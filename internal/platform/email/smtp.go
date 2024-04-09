package email

import (
	"os"

	"github.com/go-mail/mail"
	"github.com/golang/glog"
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

// Send is the method that send the email through smtp server.
func (s *SMTPService) Send(destination, name string) error {
	newMessage := mail.NewMessage()
	newMessage.SetHeader("From", s.from)
	newMessage.SetHeader("To", destination)
	newMessage.SetAddressHeader("Cc", destination, name)
	newMessage.SetHeader("Subject", Subject)
	tmpl, err := s.getTemplate()
	if err != nil {
		return err
	}
	newMessage.SetBody("text/html", tmpl)
	d := mail.NewDialer(s.host, s.port, s.username, s.password)
	if err := d.DialAndSend(newMessage); err != nil {
		glog.Errorf("error sending the email [%s]", err)

		return err
	}

	return nil
}
