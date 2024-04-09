package email

import (
	"bytes"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
	"os"
	"text/template"

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

// getTemplate reads the content of the template.
func (s *SMTPService) getTemplate() (string, error) {
	data, err := os.ReadFile(s.templatePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Send is the method that send the email through smtp server.
func (s *SMTPService) Send(destination, name string, balance mooc.Balance) error {
	newMessage := mail.NewMessage()
	newMessage.SetHeader("From", s.from)
	newMessage.SetHeader("To", destination)
	newMessage.SetAddressHeader("Cc", destination, name)
	newMessage.SetHeader("Subject", Subject)
	tmpl, err := s.getTemplate()
	if err != nil {
		return err
	}
	t := template.Must(template.New("email").Parse(tmpl))
	realValues := map[string]interface{}{
		"Name":                name,
		"Balance":             balance.Total,
		"AverageDebitMount":   balance.AverageDebitAmount,
		"AverageCreditAmount": balance.AverageCreditAmount,
		"TransactionByMonth":  balance.TransactionByMonth,
	}
	// buffer for new replaced string
	var strBuffer bytes.Buffer
	// replace the values
	err = t.Execute(&strBuffer, realValues)
	if err != nil {
		return err
	}

	newMessage.SetBody("text/html", strBuffer.String())
	d := mail.NewDialer(s.host, s.port, s.username, s.password)
	if err := d.DialAndSend(newMessage); err != nil {
		glog.Errorf("error sending the email [%s]", err)

		return err
	}

	return nil
}
