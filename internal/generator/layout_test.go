package generator

import (
	"os"
	"path/filepath"
	"testing"
)

// nestedLayout builds a templates dir in the current (base + types) shape.
func nestedLayout(t *testing.T, types ...string) Layout {
	t.Helper()

	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, "base"), 0o755); err != nil {
		t.Fatalf("mkdir base: %v", err)
	}

	for _, name := range types {

		dir := filepath.Join(root, "types", name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("mkdir type %s: %v", name, err)
		}

		manifest := `{"description": "` + name + ` desc"}`
		if err := os.WriteFile(filepath.Join(dir, "type.json"), []byte(manifest), 0o644); err != nil {
			t.Fatalf("write type.json: %v", err)
		}
	}

	return ResolveLayout(root)
}

func TestResolveLayoutNested(t *testing.T) {

	l := nestedLayout(t, "general", "casino")

	if l.Legacy {
		t.Error("nested layout must not be reported as legacy")
	}
	if filepath.Base(l.BaseDir) != "base" {
		t.Errorf("BaseDir = %q, want it to end in /base", l.BaseDir)
	}
	if l.TypesDir == "" {
		t.Error("TypesDir must be set when types/ exists")
	}
}

func TestResolveLayoutLegacyFlat(t *testing.T) {

	// a flat tree, as shipped before --type existed
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "main.go.tmpl"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	l := ResolveLayout(root)

	if !l.Legacy {
		t.Error("flat layout must be reported as legacy")
	}
	if l.BaseDir != root {
		t.Errorf("BaseDir = %q, want %q", l.BaseDir, root)
	}
	if l.TypesDir != "" {
		t.Errorf("TypesDir = %q, want empty in legacy mode", l.TypesDir)
	}
}

// An installer that copies over an old tree leaves both shapes present.
// base/ must win so the stale flat files are never walked.
func TestResolveLayoutMergedUpgradePrefersNested(t *testing.T) {

	root := t.TempDir()

	if err := os.WriteFile(filepath.Join(root, "main.go.tmpl"), []byte("stale"), 0o644); err != nil {
		t.Fatalf("write stale: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "base"), 0o755); err != nil {
		t.Fatalf("mkdir base: %v", err)
	}

	l := ResolveLayout(root)

	if l.Legacy {
		t.Error("merged tree with base/ present must not be legacy")
	}
	if filepath.Base(l.BaseDir) != "base" {
		t.Errorf("BaseDir = %q, want the nested base", l.BaseDir)
	}
}

func TestTypesListing(t *testing.T) {

	l := nestedLayout(t, "payment", "casino", "general")

	got := l.Types()
	if len(got) != 3 {
		t.Fatalf("got %d types, want 3", len(got))
	}

	// sorted by name
	want := []string{"casino", "general", "payment"}
	for i, w := range want {
		if got[i].Name != w {
			t.Errorf("types[%d] = %q, want %q", i, got[i].Name, w)
		}
		if got[i].Description != w+" desc" {
			t.Errorf("types[%d] description = %q, want it read from type.json", i, got[i].Description)
		}
	}
}

func TestTypesSkipsHiddenAndUnderscored(t *testing.T) {

	l := nestedLayout(t, "casino")

	for _, name := range []string{".hidden", "_scratch"} {
		if err := os.MkdirAll(filepath.Join(l.TypesDir, name), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", name, err)
		}
	}

	for _, got := range l.Types() {
		if got.Name == ".hidden" || got.Name == "_scratch" {
			t.Errorf("type listing must skip %q", got.Name)
		}
	}
}

func TestTypesEmptyInLegacyMode(t *testing.T) {

	l := ResolveLayout(t.TempDir())

	if len(l.Types()) != 0 {
		t.Error("legacy layout must report no types")
	}
}

func TestResolveTypeKnown(t *testing.T) {

	l := nestedLayout(t, "general", "casino", "payment")

	for _, name := range []string{"casino", "payment", "general"} {

		canonical, overlay, err := l.ResolveType(name)
		if err != nil {
			t.Errorf("ResolveType(%q) errored: %v", name, err)
			continue
		}
		if canonical != name {
			t.Errorf("ResolveType(%q) canonical = %q", name, canonical)
		}
		if overlay == "" {
			t.Errorf("ResolveType(%q) returned no overlay dir", name)
		}
	}
}

// The aliases must land on the general overlay. Falling through to a base-only
// generation would omit app/router/router.go and the service would not compile.
func TestResolveTypeAliasesUseGeneralOverlay(t *testing.T) {

	l := nestedLayout(t, "general", "casino")

	for _, alias := range []string{"", "none", "base", "GENERAL", "  general  "} {

		canonical, overlay, err := l.ResolveType(alias)
		if err != nil {
			t.Errorf("ResolveType(%q) errored: %v", alias, err)
			continue
		}
		if canonical != GeneralType {
			t.Errorf("ResolveType(%q) canonical = %q, want %q", alias, canonical, GeneralType)
		}
		if overlay == "" {
			t.Errorf("ResolveType(%q) must resolve to the general overlay, got none", alias)
		}
	}
}

func TestResolveTypeIsCaseAndSpaceInsensitive(t *testing.T) {

	l := nestedLayout(t, "casino")

	canonical, overlay, err := l.ResolveType("  CASINO ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if canonical != "casino" || overlay == "" {
		t.Errorf("got (%q, %q), want casino with an overlay", canonical, overlay)
	}
}

func TestResolveTypeUnknownErrorsAndLists(t *testing.T) {

	l := nestedLayout(t, "general", "casino", "payment")

	_, _, err := l.ResolveType("crypto")
	if err == nil {
		t.Fatal("unknown type must error")
	}

	for _, want := range []string{"crypto", "casino", "payment"} {
		if !contains(err.Error(), want) {
			t.Errorf("error message must mention %q, got:\n%s", want, err)
		}
	}
}

func TestResolveTypeUnknownInLegacyModeSuggestsReinstall(t *testing.T) {

	l := ResolveLayout(t.TempDir())

	_, _, err := l.ResolveType("casino")
	if err == nil {
		t.Fatal("a typed request against legacy templates must error")
	}
	if !contains(err.Error(), "legacy") {
		t.Errorf("legacy error should explain the situation, got:\n%s", err)
	}
}

func TestResolveTypeGeneralInLegacyModeStillWorks(t *testing.T) {

	l := ResolveLayout(t.TempDir())

	canonical, overlay, err := l.ResolveType("")
	if err != nil {
		t.Fatalf("untyped generation must still work on legacy templates: %v", err)
	}
	if canonical != GeneralType {
		t.Errorf("canonical = %q, want %q", canonical, GeneralType)
	}
	if overlay != "" {
		t.Errorf("legacy mode has no overlays, got %q", overlay)
	}
}

func contains(haystack, needle string) bool {

	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
