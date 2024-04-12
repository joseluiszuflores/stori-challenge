package email

import (
	"bytes"
	"os"
	"text/template"

	mooc "github.com/joseluiszuflores/stori-challenge/internal"

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

//nolint:lll
func NewSMTPService(host string, port int, username string, password string, from string, templatePath string) *SMTPService {
	return &SMTPService{host: host, port: port, username: username, password: password, from: from, templatePath: templatePath}
}

const Subject = "Total Balance"

//nolint:lll
const templateAux = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD XHTML 1.0 Transitional //EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
	Libertad<br />Tranquilidad <br />Stori
	Dear {{.Name}},
	<ol>
		<li style="line-height: 22.4px;">Total balance is {{.Balance}}</li>
		{{ range  $index, $element := .TransactionByMonth}}
		  <li style="line-height: 22.4px;">Number of transactions in {{$index}}: {{$element}}</li>
		{{end}}
		 {{ range  $index, $element := .AverageByMonth}}
		  <li style="line-height: 22.4px;">Average debit amount in {{$index}}: {{$element.AverageDebitAmount}}</li>
		  <li style="line-height: 22.4px;">Average credit amount in {{$index}} : {{$element.AverageCreditAmount}}</li>

  		 {{end}}

	</ol>
`

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
		glog.Error("error getting the template for email - ", err)
		tmpl = templateAux
	}
	tmplate := template.Must(template.New("email").Parse(tmpl))
	realValues := map[string]interface{}{
		"Name":                name,
		"Balance":             balance.Total,
		"AverageDebitMount":   balance.AverageDebitAmount,
		"AverageCreditAmount": balance.AverageCreditAmount,
		"TransactionByMonth":  balance.TransactionByMonth,
		"AverageByMonth":      balance.AverageByMonth,
	}
	// buffer for new replaced string
	var strBuffer bytes.Buffer
	// replace the values
	err = tmplate.Execute(&strBuffer, realValues)
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
