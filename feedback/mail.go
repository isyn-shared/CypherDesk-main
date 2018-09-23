package feedback

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

type MailRequest struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewMailRequest(to []string, subject string) *MailRequest {
	return &MailRequest{
		to:      to,
		subject: subject,
	}
}

func (r *MailRequest) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *MailRequest) sendMail() bool {
	c := new(config)
	c.Read()
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", c.Server, c.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", c.Email, c.Password, c.Server), c.Email, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *MailRequest) Send(templateName string, items interface{}) bool {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	if ok := r.sendMail(); ok {
		return true
	}
	return false
}
