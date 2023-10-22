package cmd

import (
	"github.com/l2dy/song-collection-scripts/actions"
	"github.com/spf13/cobra"
)

// normalizeFolderByRefCmd represents the normalizeFolderByRef command
var normalizeFolderByRefCmd = &cobra.Command{
	Use:   "normalizeFolderByRef",
	Short: "Find matching files by referencing srcDir and normalize the folder names",
	Long: `Rename folders in dstDir whose name is a different byte representation
of the same canonical form, when compared with srcDir.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return actions.NormalizeFolderByRef(srcDir, dstDir, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(normalizeFolderByRefCmd)

	// Here you will define your flags and configuration settings.
}
