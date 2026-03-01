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
	content := Build("Boss Hari", checklist, "https://example.com/logo.svg", "https://example.com/note.svg")

	if !strings.Contains(content.Plain, "Update work tracker") {
		t.Fatal("plain content should include checklist item")
	}
	if !strings.Contains(content.HTML, "<li") {
		t.Fatal("html content should include list item markup")
	}
	if !strings.Contains(content.HTML, "https://example.com/logo.svg") {
		t.Fatal("html content should include logo svg URL")
	}
	if !strings.Contains(content.HTML, "https://example.com/note.svg") {
		t.Fatal("html content should include note svg URL")
	}
}

func TestBuildFallsBackWhenChecklistHasNoBullets(t *testing.T) {
	content := Build("Boss Hari", "hello", "", "")

	if !strings.Contains(content.Plain, "Update today's work notes.") {
		t.Fatal("plain content should include fallback item")
	}
	if !strings.Contains(content.HTML, "Set tomorrow&#39;s top priorities.") {
		t.Fatal("html content should include fallback item")
	}
}
