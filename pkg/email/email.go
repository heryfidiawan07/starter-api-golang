package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"
)

type Mailer struct {
	Host     string
	Port     int
	User     string
	Pass     string
	From     string
	FromName string
}

func New(host string, port int, user, pass, from, fromName string) *Mailer {
	return &Mailer{
		Host:     host,
		Port:     port,
		User:     user,
		Pass:     pass,
		From:     from,
		FromName: fromName,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)

	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		fmt.Sprintf("From: %s <%s>\r\n", m.FromName, m.From) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" + body

	addr := m.Host + ":" + strconv.Itoa(m.Port)
	return smtp.SendMail(addr, auth, m.From, []string{to}, []byte(msg))
}

func (m *Mailer) SendVerification(to, name, verifyURL string) error {
	tmpl := `<!DOCTYPE html>
<html><body>
<h2>Hi {{.Name}},</h2>
<p>Please verify your email address by clicking the link below:</p>
<p><a href="{{.URL}}">Verify Email</a></p>
<p>This link will expire in 24 hours.</p>
<p>If you didn't register, ignore this email.</p>
</body></html>`

	t, err := template.New("verify").Parse(tmpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]string{"Name": name, "URL": verifyURL}); err != nil {
		return err
	}
	return m.Send(to, "Verify Your Email", buf.String())
}

func (m *Mailer) SendPasswordReset(to, name, resetURL string) error {
	tmpl := `<!DOCTYPE html>
<html><body>
<h2>Hi {{.Name}},</h2>
<p>You requested a password reset. Click the link below to reset your password:</p>
<p><a href="{{.URL}}">Reset Password</a></p>
<p>This link will expire in {{.Expire}} minutes.</p>
<p>If you didn't request this, ignore this email.</p>
</body></html>`

	t, err := template.New("reset").Parse(tmpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]string{"Name": name, "URL": resetURL, "Expire": "60"}); err != nil {
		return err
	}
	return m.Send(to, "Reset Your Password", buf.String())
}
