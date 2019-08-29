package mail

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/philtyl/anonymoe/models"
	"github.com/philtyl/anonymoe/pkg/setting"
	log "gopkg.in/clog.v1"
)

// The Backend implements SMTP server methods.
type Backend struct{}

// Login handles a login command with username and password.
func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return nil, smtp.ErrAuthUnsupported
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{Item: &models.RawMailItem{}}, nil
}

// A Session is returned after successful login.
type Session struct {
	Item *models.RawMailItem
}

func (s *Session) Mail(from string) error {
	s.Item.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.Item.Recipient = append(s.Item.Recipient, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	s.Item.Data = string(b)
	s.Item.Complete = true
	return nil
}

func (s *Session) Reset() {
	if s.Item.Complete {
		mail, recipients, err := models.CreateMail(s.Item)
		if err == nil {
			log.Info("Mail Received: %+v\nSent to: %+v", mail, recipients)
		} else {
			log.Info("Error Finalizing Mail Item: %v", err)
		}
	}
	s.Item = new(models.RawMailItem)
}

func (s *Session) Logout() error {
	return nil
}

func NewSMTPServer() {
	be := &Backend{}

	s := smtp.NewServer(be)
	s.Addr = ":" + setting.MailPort
	s.Domain = setting.AppDomain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AuthDisabled = true

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(2, "Failed to start mail server: %+v", err)
	}
}
