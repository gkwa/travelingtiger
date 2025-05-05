package renderer

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"text/template"

	"github.com/atotto/clipboard"
	"github.com/gkwa/travelingtiger/app/templates"
)

// Renderer handles template rendering operations
type Renderer struct {
	templateFS fs.FS
	templates  *template.Template
}

// NewRenderer creates a new renderer instance
func NewRenderer() (*Renderer, error) {
	// Parse all templates at initialization time
	tmpl := template.New("")
	tmpl, err := tmpl.ParseFS(templates.TemplateFS, "*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}
	return &Renderer{
		templateFS: templates.TemplateFS,
		templates:  tmpl,
	}, nil
}

// Render processes the named template and outputs it based on the outfile parameter
func (r *Renderer) Render(templateName, outFile string) error {
	// Check if template exists
	templatePath := fmt.Sprintf("%s.tmpl", templateName)
	_, err := fs.Stat(r.templateFS, templatePath)
	if err != nil {
		return r.handleTemplateNotFound(templateName)
	}

	// Clear buffer
	var buf bytes.Buffer

	// Execute the requested template instead of always using "page"
	err = r.templates.ExecuteTemplate(&buf, templateName, nil)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Get rendered content and clean it up
	content := cleanOutput(buf.String())

	// Output the result based on outFile parameter
	return r.handleOutput(content, outFile)
}

// cleanOutput removes excess whitespace and trailing characters from the template output
func cleanOutput(content string) string {
	// Replace Windows-style line endings with Unix-style
	content = strings.ReplaceAll(content, "\r\n", "\n")

	// Normalize spacing
	// Replace consecutive spaces with a single space
	for strings.Contains(content, "  ") {
		content = strings.ReplaceAll(content, "  ", " ")
	}

	// Normalize line endings
	// Replace three or more consecutive newlines with two newlines
	for strings.Contains(content, "\n\n\n") {
		content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
	}

	// Remove any trailing slashes, spaces, tabs at the end
	content = strings.TrimRight(content, "\t\n\r ")

	// Trim all leading/trailing whitespace
	content = strings.TrimSpace(content) + "\n"

	return content
}

// handleTemplateNotFound generates an error message with available templates
func (r *Renderer) handleTemplateNotFound(templateName string) error {
	// List available templates for error message
	templates, err := r.listTemplates()
	if err != nil {
		return fmt.Errorf("template '%s' not found and could not list available templates: %w", templateName, err)
	}
	if len(templates) == 0 {
		return fmt.Errorf("template '%s' not found and no templates are available", templateName)
	}
	return fmt.Errorf("template '%s' not found. Available templates: %s", templateName, strings.Join(templates, ", "))
}

// listTemplates returns a list of available template names
func (r *Renderer) listTemplates() ([]string, error) {
	var templates []string
	err := fs.WalkDir(r.templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".tmpl") {
			// Remove the .tmpl extension
			templates = append(templates, strings.TrimSuffix(path, ".tmpl"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// handleOutput sends the rendered content to the appropriate destination
func (r *Renderer) handleOutput(content, outFile string) error {
	switch outFile {
	case "":
		// Default: send to clipboard
		return clipboard.WriteAll(content)
	case "-":
		// Send to stdout
		fmt.Print(content)
		return nil
	default:
		// Write to specified file
		return os.WriteFile(outFile, []byte(content), 0o644)
	}
}
