package cmd

import (
	"fmt"
	"os"

	"github.com/gkwa/travelingtiger/app/renderer"
	"github.com/spf13/cobra"
)

var outFile string

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render [template-name]",
	Short: "Render a template",
	Long: `Render a template whose name is provided as an argument.
The rendered template is sent to the clipboard by default.
Use --outfile to specify an output file or '-' for stdout.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]

		r, err := renderer.NewRenderer()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error initializing renderer:", err)
			os.Exit(1)
		}

		err = r.Render(templateName, outFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringVar(&outFile, "outfile", "", "output file (default is clipboard, use '-' for stdout)")
}
