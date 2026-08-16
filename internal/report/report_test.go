package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"a11y-checker/internal/rules"
)

func TestReport_TextCounts(t *testing.T) {
	rep := &Report{Issues: []rules.Issue{
		{Rule: "img-alt", Severity: rules.Serious, Message: "missing alt"},
		{Rule: "html-lang", Severity: rules.Serious, Message: "missing lang"},
	}}
	if rep.Counts()["serious"] != 2 {
		t.Fatalf("serious count = %d", rep.Counts()["serious"])
	}
	var buf bytes.Buffer
	rep.RenderText(&buf)
	if !strings.Contains(buf.String(), "total: 2") {
		t.Fatalf("text missing total: %q", buf.String())
	}
}

func TestReport_CountsBySeverity(t *testing.T) {
	rep := &Report{Issues: []rules.Issue{
		{Rule: "img-alt", Severity: rules.Serious, Message: "x"},
		{Rule: "heading-order", Severity: rules.Moderate, Message: "y"},
	}}
	c := rep.Counts()
	if c["serious"] != 1 || c["moderate"] != 1 {
		t.Fatalf("counts=%v", c)
	}
}

func TestReport_JSON(t *testing.T) {
	rep := &Report{Issues: []rules.Issue{{Rule: "img-alt", Severity: rules.Serious, Message: "x"}}}
	var buf bytes.Buffer
	if err := rep.RenderJSON(&buf); err != nil {
		t.Fatalf("json error: %v", err)
	}
	var out []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("json invalid: %v", err)
	}
	if len(out) != 1 || out[0]["rule"] != "img-alt" {
		t.Fatalf("json wrong: %+v", out)
	}
}
