package main

import (
	"bytes"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	a := newAnalyzer()
	for _, password := range []string{"Password1", "Password1", "summer2024!", "123456"} {
		a.add(password)
	}
	if a.total != 4 || len(a.passwords) != 3 {
		t.Fatalf("unexpected totals: %d/%d", a.total, len(a.passwords))
	}
	if a.lengths[9] != 2 || a.charSets["lower+upper+number"] != 2 {
		t.Fatalf("unexpected length or character stats: lengths=%v sets=%v", a.lengths, a.charSets)
	}
}

func TestReport(t *testing.T) {
	a := newAnalyzer()
	a.add("Password1")
	var output bytes.Buffer
	a.report(&output, 1, false)
	if output.String() != "Password1\n" {
		t.Fatal("report is missing expected results")
	}
}
