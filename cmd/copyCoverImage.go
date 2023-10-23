package cmd

import (
	"github.com/l2dy/song-collection-scripts/actions"
	"github.com/spf13/cobra"
)

// copyCoverImageCmd represents the copyCoverImage command
var copyCoverImageCmd = &cobra.Command{
	Use:   "copyCoverImage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return actions.CopyCoverImage(srcDir, dstDir, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(copyCoverImageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCoverImageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCoverImageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
