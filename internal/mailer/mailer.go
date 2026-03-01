package mailer

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Config holds SMTP and email routing configuration loaded from env vars.
type Config struct {
	Host string
	Port string
	User string
	Pass string
	To   string
}

// Message is the outgoing email payload.
type Message struct {
	To      string
	Subject string
	Body    string
}

// Sender defines behavior for a mail delivery implementation.
type Sender interface {
	Send(cfg Config, msg Message) error
}

// SMTPClient sends mail using the standard net/smtp package.
type SMTPClient struct{}

func NewSMTPClient() *SMTPClient {
	return &SMTPClient{}
}

func LoadConfigFromEnv() (Config, error) {
	cfg := Config{
		Host: strings.TrimSpace(os.Getenv("SMTP_HOST")),
		Port: strings.TrimSpace(os.Getenv("SMTP_PORT")),
		User: strings.TrimSpace(os.Getenv("SMTP_USER")),
		Pass: strings.TrimSpace(os.Getenv("SMTP_PASS")),
		To:   strings.TrimSpace(os.Getenv("MAIL_TO")),
	}

	missing := make([]string, 0, 5)
	if cfg.Host == "" {
		missing = append(missing, "SMTP_HOST")
	}
	if cfg.Port == "" {
		missing = append(missing, "SMTP_PORT")
	}
	if cfg.User == "" {
		missing = append(missing, "SMTP_USER")
	}
	if cfg.Pass == "" {
		missing = append(missing, "SMTP_PASS")
	}
	if cfg.To == "" {
		missing = append(missing, "MAIL_TO")
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		return Config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) validate() error {
	if c.Host == "" || c.Port == "" || c.User == "" || c.Pass == "" || c.To == "" {
		return errors.New("incomplete mailer config")
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid SMTP_PORT value %q", c.Port)
	}

	return nil
}

func (c Config) address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (m Message) validate() error {
	if strings.TrimSpace(m.To) == "" {
		return errors.New("message recipient is empty")
	}
	if strings.TrimSpace(m.Subject) == "" {
		return errors.New("message subject is empty")
	}
	if strings.TrimSpace(m.Body) == "" {
		return errors.New("message body is empty")
	}
	return nil
}

func (s *SMTPClient) Send(cfg Config, msg Message) error {
	if err := cfg.validate(); err != nil {
		return err
	}
	if err := msg.validate(); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host)
	if err := smtp.SendMail(cfg.address(), auth, cfg.User, []string{msg.To}, buildPayload(cfg, msg)); err != nil {
		return fmt.Errorf("smtp sendmail: %w", err)
	}

	return nil
}

func buildPayload(cfg Config, msg Message) []byte {
	subject := strings.NewReplacer("\r", "", "\n", "").Replace(strings.TrimSpace(msg.Subject))
	body := strings.ReplaceAll(msg.Body, "\r\n", "\n")
	body = strings.ReplaceAll(body, "\n", "\r\n")

	lines := []string{
		fmt.Sprintf("From: %s", cfg.User),
		fmt.Sprintf("To: %s", msg.To),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}
	return []byte(strings.Join(lines, "\r\n"))
}
