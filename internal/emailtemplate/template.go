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

// Build creates a friendly reminder email from checklist content.
func Build(name, checklistContent string) Content {
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
		HTML:  buildHTML(displayName, items),
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

func buildHTML(name string, items []string) string {
	var list strings.Builder
	for _, item := range items {
		list.WriteString(`<li style="margin:0 0 10px 0;">`)
		list.WriteString(html.EscapeString(item))
		list.WriteString(`</li>`)
	}

	escapedName := html.EscapeString(name)
	icon := inlineLogoSVG()

	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;padding:0;background:#f8fafc;">
  <div style="max-width:680px;margin:24px auto;padding:0 14px;font-family:Segoe UI,Helvetica,Arial,sans-serif;color:#1f2937;">
    <div style="background:#ffffff;border:1px solid #e5e7eb;border-radius:20px;overflow:hidden;box-shadow:0 8px 24px rgba(15,23,42,0.08);">
      <div style="height:14px;background:#FF9933;"></div>
      <div style="background:#ffffff;padding:14px 18px 10px 18px;text-align:center;">
        %s
        <h1 style="margin:10px 0 0 0;text-align:center;font-size:30px;line-height:1.2;color:#000080;">Indus Reminder</h1>
      </div>
      <div style="height:14px;background:#138808;"></div>
      <div style="padding:24px 24px 16px 24px;font-size:16px;line-height:1.6;">
        <p style="margin:0 0 8px 0;">Hey %s,</p>
        <p style="margin:0 0 10px 0;">Hope you are doing fine.</p>
        <p style="margin:0 0 10px 0;">I thought to remind you about a few important entries for today.</p>
        <p style="margin:0 0 14px 0;">If everything is already done, that is wonderful. If not, please finish these before day end:</p>
        <ul style="padding-left:22px;margin:0 0 14px 0;">
          %s
        </ul>
        <div style="border:1px solid #e5e7eb;border-left:6px solid #FF9933;background:#fff7ed;padding:12px 14px;border-radius:10px;margin:0 0 10px 0;">
          If already done, no problem. Just continue your good work.
        </div>
        <div style="border:1px solid #e5e7eb;border-left:6px solid #138808;background:#f0fdf4;padding:12px 14px;border-radius:10px;margin:0 0 12px 0;">
          If something is pending, please complete your work updates and daily entries.
        </div>
        <p style="margin:0 0 10px 0;color:#000080;font-weight:600;">You are doing great. Keep moving forward.</p>
        <p style="margin:0 0 0 0;">With care,<br/>Indus Reminder</p>
      </div>
    </div>
  </div>
</body>
</html>`, icon, escapedName, list.String())
}

func inlineLogoSVG() string {
	return `<svg role="img" aria-label="Indus Reminder logo" width="64" height="64" viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg" style="display:inline-block;">
  <rect x="2" y="2" width="60" height="60" rx="30" fill="#ffffff" stroke="#000080" stroke-width="2"/>
  <rect x="10" y="16" width="44" height="8" rx="4" fill="#FF9933"/>
  <rect x="10" y="40" width="44" height="8" rx="4" fill="#138808"/>
  <circle cx="32" cy="32" r="8" fill="none" stroke="#000080" stroke-width="2"/>
  <text x="32" y="36" text-anchor="middle" style="font:700 9px 'Segoe UI', Arial, sans-serif;fill:#000080;">IR</text>
</svg>`
}
