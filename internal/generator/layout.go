package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GeneralType is the canonical name for a service generated from the base
// templates alone, with no overlay applied.
const GeneralType = "general"

// generalAliases all resolve to GeneralType when passed to --type.
var generalAliases = []string{"", GeneralType, "none", "base"}

// Layout describes how a templates directory is organised. Two shapes are
// supported: the nested layout (templates/base + templates/types/<name>), and
// the legacy flat layout shipped before --type existed, where the templates
// directory itself is the base and no overlays are available.
type Layout struct {
	Root     string
	BaseDir  string
	TypesDir string
	Legacy   bool
}

// TypeInfo is a service type discovered on disk.
type TypeInfo struct {
	Name        string
	Description string
}

// ResolveLayout inspects a templates directory and reports how it is organised.
// The presence of a base/ subdirectory is what distinguishes the nested layout
// from the legacy flat one. This also handles the upgrade case where an
// installer has merged new nested templates on top of an old flat tree: base/
// wins and the stale flat files are never walked.
func ResolveLayout(templatesDir string) Layout {

	baseDir := filepath.Join(templatesDir, "base")

	info, err := os.Stat(baseDir)
	if err != nil || !info.IsDir() {

		return Layout{Root: templatesDir, BaseDir: templatesDir, Legacy: true}
	}

	layout := Layout{Root: templatesDir, BaseDir: baseDir}

	typesDir := filepath.Join(templatesDir, "types")
	if info, err := os.Stat(typesDir); err == nil && info.IsDir() {
		layout.TypesDir = typesDir
	}

	return layout
}

// Types lists the service types available in this layout, sorted by name.
// Adding a type is a matter of creating a directory under types/ — no code
// change is required here.
func (l Layout) Types() []TypeInfo {

	types := []TypeInfo{}

	if l.TypesDir == "" {
		return types
	}

	entries, err := os.ReadDir(l.TypesDir)
	if err != nil {
		return types
	}

	for _, entry := range entries {

		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}

		types = append(types, TypeInfo{
			Name:        name,
			Description: readTypeDescription(filepath.Join(l.TypesDir, name)),
		})
	}

	sort.Slice(types, func(i, j int) bool { return types[i].Name < types[j].Name })

	return types
}

// readTypeDescription reads the optional type.json manifest describing a type.
func readTypeDescription(dir string) string {

	content, err := os.ReadFile(filepath.Join(dir, "type.json"))
	if err != nil {
		return ""
	}

	var manifest struct {
		Description string `json:"description"`
	}

	if err := json.Unmarshal(content, &manifest); err != nil {
		return ""
	}

	return manifest.Description
}

// ResolveType maps a --type value to its canonical name and overlay directory.
// An empty overlay directory means base-only generation.
func (l Layout) ResolveType(name string) (string, string, error) {

	normalized := strings.ToLower(strings.TrimSpace(name))

	if l.TypesDir != "" {

		overlay := filepath.Join(l.TypesDir, normalized)
		if info, err := os.Stat(overlay); err == nil && info.IsDir() && normalized != "" {
			return normalized, overlay, nil
		}
	}

	// The aliases resolve to the general overlay when it exists. Falling through
	// to a base-only generation would omit app/router/router.go, which each type
	// overlay supplies, and the service would not compile.
	for _, alias := range generalAliases {

		if normalized != alias {
			continue
		}

		if l.TypesDir != "" {

			overlay := filepath.Join(l.TypesDir, GeneralType)
			if info, err := os.Stat(overlay); err == nil && info.IsDir() {
				return GeneralType, overlay, nil
			}
		}

		return GeneralType, "", nil
	}

	return "", "", l.unknownTypeError(normalized)
}

func (l Layout) unknownTypeError(name string) error {

	if l.Legacy {

		return fmt.Errorf(`❌ Unknown service type %q

📦 This installation ships legacy templates with no type overlays.

💡 Reinstall gomicrogen to pick up the current templates, or omit --type to
   generate a general service:
   curl -fsSL https://raw.githubusercontent.com/surahj/gomicrogen/main/install-oneline.sh | bash`, name)
	}

	var b strings.Builder

	fmt.Fprintf(&b, "❌ Unknown service type %q\n\n📦 Available types:\n", name)
	fmt.Fprintf(&b, "   • %-10s base microservice (default)\n", GeneralType)

	for _, t := range l.Types() {

		if t.Name == GeneralType {
			continue
		}

		fmt.Fprintf(&b, "   • %-10s %s\n", t.Name, t.Description)
	}

	fmt.Fprintf(&b, "\n💡 Add your own: create templates/types/<name>/ with files mirroring the\n")
	fmt.Fprintf(&b, "   base layout; a file at the same path replaces the base version.")

	return fmt.Errorf("%s", b.String())
}
