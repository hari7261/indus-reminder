package mailer

import (
	"strings"
	"testing"
)

func TestLoadConfigFromEnvMissing(t *testing.T) {
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_PORT", "")
	t.Setenv("SMTP_USER", "")
	t.Setenv("SMTP_PASS", "")
	t.Setenv("MAIL_TO", "")

	_, err := LoadConfigFromEnv()
	if err == nil {
		t.Fatal("expected error for missing env vars")
	}

	for _, key := range []string{"SMTP_HOST", "SMTP_PORT", "SMTP_USER", "SMTP_PASS", "MAIL_TO"} {
		if !strings.Contains(err.Error(), key) {
			t.Fatalf("expected error to mention %s, got: %v", key, err)
		}
	}
}

func TestLoadConfigFromEnvSuccess(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.gmail.com")
	t.Setenv("SMTP_PORT", "587")
	t.Setenv("SMTP_USER", "sender@example.com")
	t.Setenv("SMTP_PASS", "app-password")
	t.Setenv("MAIL_TO", "receiver@example.com")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.Host != "smtp.gmail.com" {
		t.Fatalf("unexpected host: %s", cfg.Host)
	}
	if cfg.Port != "587" {
		t.Fatalf("unexpected port: %s", cfg.Port)
	}
	if cfg.To != "receiver@example.com" {
		t.Fatalf("unexpected recipient: %s", cfg.To)
	}
	if cfg.FromAddress != "sender@example.com" {
		t.Fatalf("unexpected from address: %s", cfg.FromAddress)
	}
	if cfg.FromName != "Indus Software Reminder" {
		t.Fatalf("unexpected from name: %s", cfg.FromName)
	}
}

func TestBuildPayload(t *testing.T) {
	cfg := Config{
		User:        "sender@example.com",
		FromAddress: "no-reply@indus-software.in",
		FromName:    "Indus Software Reminder",
	}
	msg := Message{
		To:       "receiver@example.com",
		Subject:  "Daily Reminder",
		Body:     "line-one\nline-two",
		HTMLBody: "<p>line-one</p>",
	}

	payload := string(buildPayload(cfg, msg))

	checks := []string{
		"From: Indus Software Reminder <no-reply@indus-software.in>",
		"To: receiver@example.com",
		"Subject: Daily Reminder",
		"Content-Type: multipart/alternative",
		"line-one\r\nline-two",
		"Content-Type: text/html; charset=UTF-8",
		"<p>line-one</p>",
	}

	for _, check := range checks {
		if !strings.Contains(payload, check) {
			t.Fatalf("payload missing %q", check)
		}
	}
}

func TestBuildPayloadPlainTextOnly(t *testing.T) {
	cfg := Config{User: "sender@example.com"}
	msg := Message{
		To:      "receiver@example.com",
		Subject: "Daily Reminder",
		Body:    "line-one\nline-two",
	}

	payload := string(buildPayload(cfg, msg))

	if !strings.Contains(payload, "Content-Type: text/plain; charset=UTF-8") {
		t.Fatal("expected plain text content type")
	}
	if strings.Contains(payload, "multipart/alternative") {
		t.Fatal("did not expect multipart content type")
	}
}

func TestLoadConfigFromEnvFromOverride(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.gmail.com")
	t.Setenv("SMTP_PORT", "587")
	t.Setenv("SMTP_USER", "sender@example.com")
	t.Setenv("SMTP_PASS", "app-password")
	t.Setenv("MAIL_TO", "receiver@example.com")
	t.Setenv("MAIL_FROM", "no-reply@indus-software.in")
	t.Setenv("MAIL_FROM_NAME", "Indus Software Reminder")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if cfg.FromAddress != "no-reply@indus-software.in" {
		t.Fatalf("unexpected from address: %s", cfg.FromAddress)
	}
	if cfg.FromName != "Indus Software Reminder" {
		t.Fatalf("unexpected from name: %s", cfg.FromName)
	}
}

func TestConfigValidateRejectsInvalidPort(t *testing.T) {
	cfg := Config{
		Host: "smtp.gmail.com",
		Port: "not-a-number",
		User: "sender@example.com",
		Pass: "app-password",
		To:   "receiver@example.com",
	}

	if err := cfg.validate(); err == nil {
		t.Fatal("expected invalid port error")
	}
}
