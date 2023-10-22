package cmd

import (
	"github.com/l2dy/song-collection-scripts/actions"
	"github.com/spf13/cobra"
)

// listLosslessTracksCmd represents the listLosslessTracks command
var listLosslessTracksCmd = &cobra.Command{
	Use:   "listLosslessTracks",
	Short: "List all lossless tracks in the collection",
	Long: `List all lossless tracks in the collection.
This only scans FLAC files and excludes unsplit files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return actions.ListLosslessTracks(srcDir, dstDir)
	},
}

func init() {
	rootCmd.AddCommand(listLosslessTracksCmd)

	// Here you will define your flags and configuration settings.
}
