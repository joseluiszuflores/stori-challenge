package config

import "github.com/jinzhu/configor"

var Config = struct {
	SMTPHost          string `default:"smtp.gmail.com" required:"true" env:"SMTP_HOST"`
	SMTPPort          int    `default:"587" required:"true" env:"SMTP_PORT"`
	SMTPUsername      string `default:"joseluiszuflores@gmail.com" required:"true" env:"SMTP_USERNAME"`
	SMTPPassword      string `default:"" required:"true" env:"SMTP_PASSWORD"`
	SMTPTemplatePath  string `default:"./../template/email/index.html" required:"true" env:"SMTP_TEMPLATE_PATH_EMAIL"`
	AWSAccessKey      string `default:"" required:"true" env:"AWS_ACCESS_KEY"`
	AWSSecretAcessKey string `default:"" required:"true" env:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion         string `default:"us-east-2" required:"true" env:"AWS_REGION"`
	DevEnv            bool   `default:"true" required:"true" env:"DEV_ENV"` // is used only dev environment or local.
	AWSUrlDynamoDev   string `default:"http://localhost:8000" required:"true" env:"AWS_URL_DYNAMO_DEV"`
}{}

func Init() error {
	if err := configor.Load(&Config); err != nil {
		return err
	}
	return nil
}
