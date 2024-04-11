package config

import "github.com/jinzhu/configor"

var Config = struct {
	SMTPHost         string `default:"smtp.gmail.com" required:"true" env:"SMTP_HOST"`
	SMTPPort         int    `default:"587" required:"true" env:"SMTP_PORT"`
	SMTPUsername     string `default:"joseluiszuflores@gmail.com" required:"true" env:"SMTP_USERNAME"`
	SMTPPassword     string `default:"" required:"true" env:"SMTP_PASSWORD"`
	SMTPTemplatePath string `default:"./../template/email/index.html" required:"true" env:"SMTP_TEMPLATE_PATH_EMAIL"`
}{}

func Init() error {
	if err := configor.Load(&Config); err != nil {
		return err
	}
	return nil
}
