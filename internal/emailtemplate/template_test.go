package emailtemplate

import (
	"fmt"
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
	if !strings.Contains(content.HTML, "<svg") {
		t.Fatal("html content should include inline svg icon")
	}
	if strings.Count(content.HTML, "<svg") < 2 {
		t.Fatal("html content should include icon svg and illustration svg")
	}
	if strings.Contains(content.HTML, "<text") {
		t.Fatal("svg graphics should not use custom text nodes")
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

func TestBuildRendersChecklistListItems(t *testing.T) {
	content := Build("Boss Hari", "- [ ] One\n- [ ] Two\n- [ ] Three")

	if got := strings.Count(content.HTML, "<li"); got != 3 {
		t.Fatal(fmt.Sprintf("expected 3 list items, got %d", got))
	}
}
