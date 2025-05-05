package lister

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/gkwa/travelingtiger/app/templates"
)

// Lister handles template listing operations
type Lister struct {
	templateFS fs.FS
}

// NewLister creates a new lister instance
func NewLister() *Lister {
	return &Lister{
		templateFS: templates.TemplateFS,
	}
}

// ListTemplates returns a list of available template names
func (l *Lister) ListTemplates() ([]string, error) {
	var templateNames []string

	// Walk the embedded filesystem to find all template files
	err := fs.WalkDir(l.templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			// Remove the .tmpl extension
			templateNames = append(templateNames, strings.TrimSuffix(path, ".tmpl"))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort templates for consistent output
	sort.Strings(templateNames)

	return templateNames, nil
}

// FindSubstringMatches returns a list of template names that contain the given substring
func (l *Lister) FindSubstringMatches(substring string) ([]string, error) {
	allTemplates, err := l.ListTemplates()
	if err != nil {
		return nil, fmt.Errorf("error finding matching templates: %w", err)
	}

	var matches []string
	lowercaseSubstr := strings.ToLower(substring)

	for _, tmpl := range allTemplates {
		if strings.Contains(strings.ToLower(tmpl), lowercaseSubstr) {
			matches = append(matches, tmpl)
		}
	}

	return matches, nil
}

// FormatTemplateList formats a slice of template names as a newline-delimited list
func (l *Lister) FormatTemplateList(templates []string) string {
	var builder strings.Builder
	for _, tmpl := range templates {
		builder.WriteString("  ")
		builder.WriteString(tmpl)
		builder.WriteString("\n")
	}
	return builder.String()
}
