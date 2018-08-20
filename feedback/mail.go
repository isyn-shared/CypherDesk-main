package feedback

import (
	"CypherDesk-main/alias"
	"net/smtp"
	"strings"
)

// MailConnection contains login, pass, host and port to use smtp
type MailConnection struct {
	mail string
	pass string
	host string
	port string
}

// MailMessage descrides the mail messsage object
type MailMessage struct {
	Subject    string
	Body       string
	Recipients []string
}

// GetString returns byte array with recipients, subj, body
func (mm *MailMessage) getString() []byte {
	res := "To: "
	for _, rec := range mm.Recipients {
		res += rec + "\r\n"
	}
	res += "Subject: " + mm.Subject + "\r\n" +
		"\r\n" +
		mm.Body + "\r\n"
	return []byte(res)
}

func getMailKey() *MailConnection {
	bs := chk(alias.ReadFile("feedback/mail.key"))

	str := bs.(string)
	lp := strings.Split(str, ";")

	if len(lp) != 4 {
		panic("Неправильный формат файла mail.key!")
	}
	mc := &MailConnection{
		mail: lp[0],
		pass: lp[1],
		host: lp[2],
		port: lp[3],
	}
	return mc
}

// SendMail sends mailMessage object using smtp
func SendMail(mm *MailMessage) error {
	mc := getMailKey()
	auth := smtp.PlainAuth(
		"",
		mc.mail,
		mc.pass,
		mc.host,
	)
	err := smtp.SendMail(mc.host+":"+mc.port, auth, mc.mail, mm.Recipients, mm.getString())
	if err != nil {
		panic(err.Error())
	}
	return nil
}
