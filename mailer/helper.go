package mailer

import (
	"bytes"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

type bot interface {
	GetTemplate(templateName string) (string, error)
}

type UserEmail struct {
	FirstName string
	URL       string
	Subject   string
	AppName   string
}

func PrintTemplate(b bot, templateName string) string {
	result, err := b.GetTemplate(templateName)
	if err != nil {
		panic(err)
	}
	return result
}

func (u UserEmail) GetTemplate(templateName string) (string, error) {
	var tpl bytes.Buffer
	t, err := ParseTemplateDir("mailer/templates")
	if err != nil {
		panic(err)
	}
	u.AppName = os.Getenv("APP_NAME")
	err = t.ExecuteTemplate(&tpl, templateName, &u)
	if err != nil {
		log.Println(err)
		return "", err
	}
	result := tpl.String()
	return result, nil
}

func SendMail(from string, to string, subject string, body string, cc string) {
	m := gomail.NewMessage()
	m.SetHeader("From", from, os.Getenv("MAIL_FROM_NAME"))
	m.SetHeader("To", to)
	if cc != "" {
		m.SetHeader("Cc", cc)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.AddAlternative("text/plain", html2text.HTML2Text(body))

	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		panic(err)
	}
	mailHost := os.Getenv("MAIL_HOST")
	mailUser := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
