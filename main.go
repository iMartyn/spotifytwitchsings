package main

import (
	"fmt"
	"os"

	spotifytwitchsings "github.com/iMartyn/spotifytwitchsings/src"
	"github.com/spf13/cobra"
)

func main() {
	var playlistid string
	var albumid string

	var rootCmd = &cobra.Command{
		Use:   "spotifytwitchsings",
		Short: "spotifytwitchsings is a awesome music downloader",
		Long: `Spotifytwitchsings lets you find which songs are available on twitch sings
Pass Either album ID or Playlist ID to start comparing`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(playlistid) > 0 && len(albumid) > 0 {
				fmt.Println("Either album ID or playlist ID")
				cmd.Help()
			} else if len(albumid) > 0 {
				// Download album with the given album ID
				spotifytwitchsings.DownloadAlbum(albumid)
			} else if len(playlistid) > 0 {
				// Download playlist with the given ID
				spotifytwitchsings.DownloadPlaylist(playlistid)
			} else {
				fmt.Println("Enter valid input.")
				cmd.Help()
			}
		},
	}

	rootCmd.Flags().StringVarP(&playlistid, "playlistid", "p", "", "Album ID found on spotify")
	rootCmd.Flags().StringVarP(&albumid, "albumid", "a", "", "Album ID found on spotify")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
