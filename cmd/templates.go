package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/gkwa/travelingtiger/app/templates"
	"github.com/spf13/cobra"
)

// templatesCmd represents the templates command
var templatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"t"},
	Short:   "List all available templates",
	Long:    `List all templates that are embedded in the binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		// List available templates
		templateList, err := listTemplates()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error listing templates:", err)
			os.Exit(1)
		}

		if len(templateList) == 0 {
			fmt.Fprintln(os.Stderr, "No templates found.")
			os.Exit(1)
		}

		// Sort templates for consistent output
		sort.Strings(templateList)

		// Output template names to stdout
		fmt.Println("Available templates:")
		for _, tmpl := range templateList {
			fmt.Println("  -", tmpl)
		}
	},
}

// listTemplates returns a list of available template names
func listTemplates() ([]string, error) {
	var templateNames []string

	// Walk the embedded filesystem to find all template files
	err := fs.WalkDir(templates.TemplateFS, ".", func(path string, d fs.DirEntry, err error) error {
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

	return templateNames, nil
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
