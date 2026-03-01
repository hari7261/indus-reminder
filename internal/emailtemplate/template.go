package emailtemplate

import (
	"fmt"
	"html"
	"strings"
)

// Content carries both plain text and HTML versions of the reminder email.
type Content struct {
	Plain string
	HTML  string
}

// Build creates a friendly reminder email from checklist content and SVG links.
func Build(name, checklistContent, logoURL, noteURL string) Content {
	displayName := strings.TrimSpace(name)
	if displayName == "" {
		displayName = "Boss Hari"
	}

	items := checklistItems(checklistContent)
	if len(items) == 0 {
		items = []string{
			"Update today's work notes.",
			"Complete today's personal entries.",
			"Set tomorrow's top priorities.",
		}
	}

	return Content{
		Plain: buildPlain(displayName, items),
		HTML:  buildHTML(displayName, items, logoURL, noteURL),
	}
}

func checklistItems(checklistContent string) []string {
	lines := strings.Split(strings.ReplaceAll(checklistContent, "\r\n", "\n"), "\n")
	items := make([]string, 0, len(lines))

	prefixes := []string{
		"- [ ] ",
		"- [x] ",
		"- [X] ",
		"* [ ] ",
		"* [x] ",
		"* [X] ",
		"- ",
		"* ",
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		for _, prefix := range prefixes {
			if strings.HasPrefix(trimmed, prefix) {
				value := strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				if value != "" {
					items = append(items, value)
				}
				break
			}
		}
	}

	return items
}

func buildPlain(name string, items []string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Hey %s,\n\n", name))
	b.WriteString("Hope you are doing fine.\n")
	b.WriteString("A gentle reminder for today.\n\n")
	b.WriteString("If everything is done, amazing. Keep going.\n")
	b.WriteString("If anything is pending, please complete these before you sign off:\n\n")
	for _, item := range items {
		b.WriteString("- ")
		b.WriteString(item)
		b.WriteString("\n")
	}
	b.WriteString("\nWith care,\nIndus Reminder\n")
	return b.String()
}

func buildHTML(name string, items []string, logoURL string, noteURL string) string {
	var list strings.Builder
	for _, item := range items {
		list.WriteString(`<li style="margin:0 0 10px 0;">`)
		list.WriteString(html.EscapeString(item))
		list.WriteString(`</li>`)
	}

	escapedName := html.EscapeString(name)
	escapedLogoURL := html.EscapeString(strings.TrimSpace(logoURL))
	escapedNoteURL := html.EscapeString(strings.TrimSpace(noteURL))

	logoBlock := ""
	if escapedLogoURL != "" {
		logoBlock = fmt.Sprintf(
			`<div style="text-align:center;margin-bottom:10px;"><img src="%s" alt="Indus Reminder" style="max-width:260px;width:100%%;height:auto;border:0;display:inline-block;"/></div>`,
			escapedLogoURL,
		)
	}

	noteBlock := ""
	if escapedNoteURL != "" {
		noteBlock = fmt.Sprintf(
			`<div style="text-align:center;margin:20px 0 8px 0;"><img src="%s" alt="Daily reminder note" style="max-width:520px;width:100%%;height:auto;border:0;display:inline-block;border-radius:14px;"/></div>`,
			escapedNoteURL,
		)
	}

	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#f4f7fb;">
  <div style="max-width:680px;margin:24px auto;padding:0 14px;font-family:Segoe UI,Helvetica,Arial,sans-serif;color:#1f2937;">
    <div style="background:#ffffff;border:1px solid #e5e7eb;border-radius:20px;overflow:hidden;">
      <div style="background:linear-gradient(135deg,#0ea5e9,#22c55e);padding:18px 18px 12px 18px;">
        %s
        <h1 style="margin:0;text-align:center;font-size:30px;line-height:1.2;color:#ffffff;">Indus Reminder</h1>
      </div>
      <div style="padding:24px 24px 10px 24px;font-size:16px;line-height:1.6;">
        <p style="margin:0 0 8px 0;">Hey %s,</p>
        <p style="margin:0 0 10px 0;">Hope you are doing fine.</p>
        <p style="margin:0 0 10px 0;">I thought to remind you about a few important entries for today.</p>
        <p style="margin:0 0 14px 0;">If everything is already done, that is wonderful. If not, please finish these before day end:</p>
        <ul style="padding-left:22px;margin:0 0 6px 0;">
          %s
        </ul>
        %s
        <p style="margin:0 0 10px 0;">You are doing great. Keep moving forward.</p>
        <p style="margin:0 0 0 0;">With care,<br/>Indus Reminder</p>
      </div>
    </div>
  </div>
</body>
</html>`, logoBlock, escapedName, list.String(), noteBlock)
}
