package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hari7261/indus-reminder/internal/emailtemplate"
	"github.com/hari7261/indus-reminder/internal/mailer"
)

const (
	defaultChecklistFile = "checklist.md"
	defaultTimezone      = "Asia/Kolkata"
	defaultSubject       = "Indus Reminder: Update work trackers and personal checklist"
	defaultReminderName  = "Boss Hari"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return errors.New("indus-reminder only runs inside GitHub Actions")
	}

	tz := envOrDefault("REMINDER_TZ", defaultTimezone)
	location, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("load timezone %q: %w", tz, err)
	}

	now := time.Now().In(location)
	forceSend := strings.EqualFold(strings.TrimSpace(os.Getenv("FORCE_SEND")), "true")
	if now.Weekday() == time.Sunday && !forceSend {
		log.Printf("skip: Sunday in timezone %s", tz)
		return nil
	}
	if now.Weekday() == time.Sunday && forceSend {
		log.Printf("force_send enabled: bypassing Sunday skip in timezone %s", tz)
	}

	checklistPath := envOrDefault("CHECKLIST_FILE", defaultChecklistFile)
	checklistBody, err := os.ReadFile(checklistPath)
	if err != nil {
		return fmt.Errorf("read checklist file %q: %w", checklistPath, err)
	}
	if strings.TrimSpace(string(checklistBody)) == "" {
		return fmt.Errorf("checklist file %q is empty", checklistPath)
	}

	cfg, err := mailer.LoadConfigFromEnv()
	if err != nil {
		return err
	}

	msg := mailer.Message{
		To:      cfg.To,
		Subject: envOrDefault("MAIL_SUBJECT", defaultSubject),
	}
	templateContent := emailtemplate.Build(
		envOrDefault("REMINDER_NAME", defaultReminderName),
		string(checklistBody),
	)
	msg.Body = templateContent.Plain
	msg.HTMLBody = templateContent.HTML

	client := mailer.NewSMTPClient()
	if err := client.Send(cfg, msg); err != nil {
		return fmt.Errorf("send reminder email: %w", err)
	}

	log.Printf("success: reminder sent to %s", cfg.To)
	return nil
}

func envOrDefault(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}
