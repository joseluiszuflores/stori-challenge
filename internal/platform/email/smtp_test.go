//go:build integration
// +build integration

package email

import (
	"os"
	"testing"
)

func TestSMTPService_Send(t *testing.T) {
	type fields struct {
		host     string
		port     int
		username string
		password string
		from     string
	}
	type args struct {
		destination string
		name        string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should send the email when all data is correct",
			fields: fields{
				host:     "smtp.gmail.com",
				port:     587,
				username: "joseluiszuflores@gmail.com",
				password: os.Getenv("psswd_email"),
				from:     "joseluiszuflores@gmail.com",
			},
			args: args{
				destination: "storymockexample@gmail.com",
				name:        "Story",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SMTPService{
				host:         tt.fields.host,
				port:         tt.fields.port,
				username:     tt.fields.username,
				password:     tt.fields.password,
				from:         tt.fields.from,
				templatePath: "./../../../template/email/index.html",
			}
			if err := s.Send(tt.args.destination, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
