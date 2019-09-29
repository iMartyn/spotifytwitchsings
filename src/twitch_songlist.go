package spotifytwitchsings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const twitchCatalogUrl = "https://sings-extension.twitch.tv/v1/catalog?sortType=artist&cursor="
const cacheTTL = 3600                  // one hour
var cachedSongList []twitchSingsSong   // runtime global for simple caching
var cachedSongListUpdateTime time.Time // when to refresh cache
var cacheFileLocation string

type MatchType int

const (
	MatchNoMatch           MatchType = 0
	MatchTrackNameOnly     MatchType = 1
	MatchBothNameAndArtist MatchType = 2
	MatchNameOnlyFuzzy     MatchType = 3 // Not yet implemented
	MatchBothFuzzy         MatchType = 4 // Not yet implemented
)

type twitchSingsResponse struct {
	Cursor  string            `json:"cursor"`
	Results []twitchSingsSong `json:"results"`
}

type twitchSingsSong struct {
	Uuid           string   `json:"uuid"`
	Name           string   `json:"name"`
	Artist         string   `json:"artist"`
	Genres         []string `json:"genres"`
	Origin         string   `json:"origin"`
	Year           int      `json:"year"`
	FirstPublished string   `json:"firstPublished"`
	Languages      []string `json:"languages"`
	HasLeadVocals  bool     `json:"hasLeadVocals,omitempty"`
}

func getTwitchResponse(cursor string) (twitchSingsResponse, error) {
	response := twitchSingsResponse{}

	twitchClient := http.Client{Timeout: time.Second * 2}
	req, err := http.NewRequest(http.MethodGet, twitchCatalogUrl+cursor, nil)
	if err != nil {
		return response, err
	}

	req.Header.Set("User-Agent", "spotifytwitchsings 0.1")
	res, getErr := twitchClient.Do(req)
	if getErr != nil {
		return response, getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return response, readErr
	}

	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		return response, jsonErr
	}
	return response, nil
}

func TwitchGetSongs() (songlist []twitchSingsSong, err error) {
	var allSongs []twitchSingsSong
	songCount := 100
	cursor := ""
	for songCount > 0 {
		response, err := getTwitchResponse(cursor)
		if err != nil {
			return songlist, err
		}
		allSongs = append(allSongs, response.Results...)
		songCount = len(response.Results)
		cursor = response.Cursor
	}
	return allSongs, nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func getCacheFileLocation() string {
	finalPath := ""
	if len(os.Getenv("HOME")) == 0 {
		finalPath = "/tmp/"
	} else {
		result, err := exists(os.Getenv("HOME")+"/.cache")
		if err != nil {
			result = false
		}
		if result {
			finalPath = os.Getenv("HOME")+"/.cache"
		} else {
			finalPath = "/tmp/"
		}
	}
	return finalPath
}

func cacheTwitchSongsToFile() error {
	bytes, err := json.Marshal(cachedSongList)
	if err != nil {
		fmt.Print("Error serialising the song list : ")
		fmt.Println(err)
		return err
	}
	err = ioutil.WriteFile(getCacheFileLocation()+"/twitchsingslist.json", bytes, 0644)
	if err != nil {
		fmt.Print("Error writing the song list : ")
		fmt.Println(err)
		return err
	}
	return err
}

func getCachedSongsFromFile() (songlist []twitchSingsSong, err error) {
	bytes, err := ioutil.ReadFile(getCacheFileLocation() + "/twitchsingslist.json")
	if err != nil {
		fmt.Print("Error reading the song list cache : ")
		fmt.Println(err)
		return songlist, err
	}
	jsonErr := json.Unmarshal(bytes, &songlist)
	if jsonErr != nil {
		fmt.Print("Error deserialising the song list : ")
		fmt.Println(jsonErr)
		return songlist, jsonErr
	}
	f, err := os.OpenFile(getCacheFileLocation() + "/twitchsingslist.json", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Print("Error opening the song list cache : ")
		fmt.Println(err)
		return songlist, err
	}
	defer f.Close()
	fStats, err := f.Stat()
	if err != nil {
		fmt.Print("Error stat()ing the song list cache : ")
		fmt.Println(err)
		return songlist, err
	}
	cachedSongListUpdateTime = fStats.ModTime().Add(time.Second * cacheTTL)
	return songlist, nil
}

func CachedTwitchGetSongs(fromFile bool) (songlist []twitchSingsSong, err error) {
	if fromFile {
		cachedSongList, _ = getCachedSongsFromFile()
	}
	if len(cachedSongList) < 0 || time.Now().After(cachedSongListUpdateTime) {
		newList, err := TwitchGetSongs()
		if err != nil {
			return songlist, err
		} else {
			cachedSongList = newList
			_ = cacheTwitchSongsToFile()
			cachedSongListUpdateTime = time.Now().Add(time.Second * cacheTTL)
			return newList, nil
		}
	} else {
		return cachedSongList, nil
	}
}

func SpotifyListContains(trackName string, artistName []string) (MatchType, twitchSingsSong) {
	artistMatches := false
	for _, val := range cachedSongList {
		if val.Name == trackName {
			for _, artist := range artistName {
				if val.Artist == artist {
					artistMatches = true
				}
			}
			if artistMatches {
				return MatchBothNameAndArtist, val
			}
			return MatchTrackNameOnly, val
		}
	}
	return MatchNoMatch, twitchSingsSong{}
}
