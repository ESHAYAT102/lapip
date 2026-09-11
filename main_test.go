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
	a.report(&output, 1, false, false)
	if output.String() != "Password1\n" {
		t.Fatal("report is missing expected results")
	}
}

func TestNumberVariants(t *testing.T) {
	a := newAnalyzer()
	a.add("P")
	var output bytes.Buffer
	a.report(&output, 1, false, true)
	result := output.String()
	for _, password := range []string{"P123\n", "P1234\n", "P.789\n", "P.6789\n", "P@000\n", "P#1234\n", "P$9999\n", "P!123\n", "P%123\n", "P&123\n", "P*123\n", "P_123\n", "P-123\n", "P+123\n", "P=123\n"} {
		if !bytes.Contains([]byte(result), []byte(password)) {
			t.Fatalf("missing %q", password)
		}
	}
	for _, password := range []string{"p\n", "p@123\n", "p#1234\n"} {
		if !bytes.Contains(output.Bytes(), []byte(password)) {
			t.Fatalf("missing %q", password)
		}
	}
	if lines := bytes.Count(output.Bytes(), []byte{'\n'}); lines != 286002 {
		t.Fatalf("got %d lines, want 286002", lines)
	}
}

func TestEachVariant(t *testing.T) {
	got := make(map[string]bool)
	eachVariant("PASS", func(word string) {
		if got[word] {
			t.Fatalf("duplicate variant %q", word)
		}
		got[word] = true
	})
	if len(got) != 54 {
		t.Fatalf("got %d variants, want 54", len(got))
	}
	for _, word := range []string{"P@SS", "P@5S", "P@S5", "P@55", "PA5S", "PAS5", "PA55", "PASS", "p@ss", "p@5s", "p@s5", "p@55", "pa5s", "pas5", "pa55", "pass", "pA5s"} {
		if !got[word] {
			t.Errorf("missing %q", word)
		}
	}
	count := 0
	eachVariant("123!", func(word string) {
		count++
		if word != "123!" {
			t.Errorf("unexpected variant %q", word)
		}
	})
	if count != 1 {
		t.Fatalf("unchanged input emitted %d times", count)
	}
}
