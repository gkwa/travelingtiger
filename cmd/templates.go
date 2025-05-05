package cmd

import (
	"fmt"
	"os"

	"github.com/gkwa/travelingtiger/app/templates/lister"
	"github.com/spf13/cobra"
)

// templatesCmd represents the templates command
var templatesCmd = &cobra.Command{
	Use:     "templates",
	Aliases: []string{"t"},
	Short:   "List all available templates",
	Long:    `List all templates that are embedded in the binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a new lister
		l := lister.NewLister()

		// List available templates
		templateList, err := l.ListTemplates()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error listing templates:", err)
			os.Exit(1)
		}

		if len(templateList) == 0 {
			fmt.Fprintln(os.Stderr, "No templates found.")
			os.Exit(1)
		}

		// Output template names to stdout
		fmt.Println("Available templates:")
		for _, tmpl := range templateList {
			fmt.Println("  -", tmpl)
		}
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
