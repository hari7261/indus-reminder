package emailtemplate

import (
	"strings"
	"testing"
)

func TestBuildUsesChecklistItems(t *testing.T) {
	checklist := `
# Daily

- [ ] Update work tracker
- [ ] Add personal entry
`
	content := Build("Boss Hari", checklist)

	if !strings.Contains(content.Plain, "Update work tracker") {
		t.Fatal("plain content should include checklist item")
	}
	if !strings.Contains(content.HTML, "<li") {
		t.Fatal("html content should include list item markup")
	}
	if !strings.Contains(content.HTML, "Indus Reminder") {
		t.Fatal("html content should include heading")
	}
	if !strings.Contains(content.HTML, "#FF9933") || !strings.Contains(content.HTML, "#138808") {
		t.Fatal("html content should use tricolour palette")
	}
}

func TestBuildFallsBackWhenChecklistHasNoBullets(t *testing.T) {
	content := Build("Boss Hari", "hello")

	if !strings.Contains(content.Plain, "Update today's work notes.") {
		t.Fatal("plain content should include fallback item")
	}
	if !strings.Contains(content.HTML, "Set tomorrow&#39;s top priorities.") {
		t.Fatal("html content should include fallback item")
	}
}
