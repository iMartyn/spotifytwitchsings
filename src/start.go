package spotifytwitchsings

import (
	"fmt"

	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
	"github.com/zmb3/spotify"
)

func SpotifyGetPlaylistTracks(pid string) []spotify.FullTrack {
	user := InitAuth()
	cli := UserData{
		UserClient: user,
	}
	playlistID := spotify.ID(pid)
	trackListJSON, _ := cli.UserClient.GetPlaylistTracks(playlistID)
	for _, val := range trackListJSON.Tracks {
		cli.TrackList = append(cli.TrackList, val.Track)
	}
	return cli.TrackList
}

func SpotifyArtistsAsString(artists []spotify.SimpleArtist) string {
	ret := ""
	for _, artist := range(artists) {
		ret += artist.Name
	}
	return ret
}

// DownloadPlaylist Start initializes complete program
func DownloadPlaylist(pid string) {
	user := InitAuth()
	cli := UserData{
		UserClient: user,
	}
	playlistID := spotify.ID(pid)
	trackListJSON, _ := cli.UserClient.GetPlaylistTracks(playlistID)
	for _, val := range trackListJSON.Tracks {
		cli.TrackList = append(cli.TrackList, val.Track)
	}
	CompareTrackList(cli)
}

// DownloadAlbum Download album according to
func DownloadAlbum(aid string) {
	user := InitAuth()
	cli := UserData{
		UserClient: user,
	}
	albumid := spotify.ID(aid)
	album, _ := user.GetAlbum(albumid)
	for _, val := range album.Tracks.Tracks {
		cli.TrackList = append(cli.TrackList, spotify.FullTrack{
			SimpleTrack: val,
			Album:       album.SimpleAlbum,
		})
	}
	CompareTrackList(cli)
}

// CompareTrackList Start downloading given list of tracks
func CompareTrackList(cli UserData) {
	if len(cachedSongList) <= 0 {
		fmt.Println("No twitch sings cached yet, reading from disk or remote...")
		CachedTwitchGetSongs(true)
	}
	fmt.Println("Found", len(cli.TrackList), "tracks")
	fmt.Println("Searching and downloading tracks")
	uiprogress.Start()
	bar := uiprogress.AddBar(len(cli.TrackList))

	bar.AppendCompleted()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		if b.Current() == len(cli.TrackList) {
			return "   🔍 " + strutil.Resize("Search complete", 30)
		}
		return "   🔍 " + strutil.Resize(cli.TrackList[b.Current()].Name, 30)
	})
	uiprogress.Stop()
	fmt.Println("Tracks : ")
	for _, val := range cli.TrackList {
		artists := []string{}
		for _, artistforname := range val.Artists {
			artists = append(artists, artistforname.Name)
		}
		matchKind, twitchSingsSong := SpotifyListContains(val.Name, artists)
		if matchKind != MatchNoMatch {
			fmt.Print(twitchSingsSong.Name + " by " + twitchSingsSong.Artist + " matches ")
			if matchKind == MatchBothNameAndArtist {
				fmt.Println("both track name and artist")
			} else {
				fmt.Println("just track name")
			}
		}
	}
}
