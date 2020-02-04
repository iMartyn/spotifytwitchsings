package spotifytwitchsings

import (
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"time"
	"strconv"
)

func sendHTMLPreambleAndHead(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type","text/html")
	fmt.Fprint(writer,
`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="Check Spotify playlists against the twitch music database">
    <meta name="author" content="Martyn Ranyard">
    <title>SpotifyTwitchSings</title>

    <!-- Bootstrap core CSS -->
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" rel="stylesheet" crossorigin="anonymous">
    <!-- Custom styles for this template -->
    <link href="/cover.css" rel="stylesheet">

    <style>
      .bd-placeholder-img {
        font-size: 1.125rem;
        text-anchor: middle;
        -webkit-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
        user-select: none;
      }

      @media (min-width: 768px) {
        .bd-placeholder-img-lg {
          font-size: 3.5rem;
        }
			}
			
			li.match-nomatch{
				background-color: #1e2122;
			}
			li.match-matchtrack{
				background-color: #E9B000;
			}
			li.match-fullmatch{
				background-color: #008F95;
			}
			li.match-matchtrackfuzzt{
				background-color: darkgray;
			}
			li.match-fullmatchfuzzy{
				background-color: darkgray;
			}
			a{
				text-decoration-line: underline;
			}

    </style>
  </head>
  `);
}

func sendBodyAndHeader(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type","text/html")
	fmt.Fprint(writer,
`<body class="text-center">
<div class="cover-container d-flex w-100 h-100 p-3 mx-auto flex-column">
<header class="masthead mb-auto">
<div class="inner">
  <h3 class="masthead-brand">SpotifyTwitchSings</h3>
  <nav class="nav nav-masthead justify-content-center">
	<a class="nav-link active" href="/">Home</a>
  </nav>
</div>
</header>
  `);
}

func sendAllTheRest(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type","text/html")
	fmt.Fprint(writer,
`<footer class="mastfoot mt-auto">
<div class="inner">
  <p>Cover template for <a href="https://getbootstrap.com/">Bootstrap</a>, by <a href="https://twitter.com/mdo">@mdo</a>.</p>
</div>
</footer>
</div>
<script>
function directToResults() {
	var url = document.createElement('a');
	url.setAttribute("href", window.location.href);
	if ((url.port != 80) && (url.port != 443)) {
		customPort = ":"+url.port
	} else {
		customPort = ""
	}
	if (document.getElementById("spotifyid").value.indexOf("album") > 0) {
		searchMode = "album"
	} else if (document.getElementById("spotifyid").value.indexOf("playlist") > 0) {
		searchMode = "playlist"
	} else {
		searchMode = document.getElementById("mode").value
	}
	pattern = /https:\/\/open\.spotify\.com\/(album|playlist)\/([0-9a-zA-Z]{22})/
	if (pattern.test(document.getElementById("spotifyid").value)) {
		parts = pattern.exec(document.getElementById("spotifyid").value)
		justTheId = parts[2]
	} else {
		justTheId = document.getElementById("spotifyid").value
	}
	var destination = url.protocol + "//" + url.hostname + customPort + "/" + searchMode + "/" + justTheId
	window.location.href = destination
}

function toggleUnfound() {
	var unmatched = document.getElementsByClassName('match-nomatch'), i;
	if (document.getElementById("showhidebutton").getAttribute("tracksHidden") != "true") {
		document.getElementById("showhidebutton").setAttribute("tracksHidden","true")
		for (i = 0; i < unmatched.length; i += 1) {
				unmatched[i].style.display = 'none';
		}
	} else {
		document.getElementById("showhidebutton").setAttribute("tracksHidden","false")
		for (i = 0; i < unmatched.length; i += 1) {
				unmatched[i].style.display = 'list-item';
		}
	}
}
</script>
</body>
</html>
  `);
}

func HomeHandler(response http.ResponseWriter, request *http.Request) {
	sendHTMLPreambleAndHead(response)
	sendBodyAndHeader(response)
	fmt.Fprint(response,
`  <main role="main" class="inner cover">
	<h1 class="cover-heading">Gimme your ids!</h1>
	<label for="spotifyid">ID of playlist or album :&nbsp;</label><input type="text" name="spotifyid" id="spotifyid">
	<select name="mode" id="mode">
	  <option value="playlist">Search playlist for twitch sings songs</option>
	  <option value="album">Search album for twitch sings songs</option>
	</select>
	<p/>
	<p class="lead">
      <a href="#" class="btn btn-lg btn-secondary" onClick="javascript:directToResults()">Search</a>
    </p>
    <p class="lead">Spotify IDs look like this : 37i9dQZF1DX4UtSsGT1Sbe - 22 characters and can be got by clicking share and "Copy (Playlist|Album) Link".</p>
		<p>Shameless self-promotion : Follow me on twitch - <a href="https://www.twitch.tv/iMartynOnTwitch">iMartynOnTwitch</a>, oddly enough, I do a lot of twitchsings!</p>
  </main>
  `)
  sendAllTheRest(response)
}

func NotFoundHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(404)
	sendHTMLPreambleAndHead(response)
	sendBodyAndHeader(response)
	fmt.Fprintf(response,
`<main role="main" class="inner cover">
	<h1 class="cover-heading">Ooops!</h1>
	<p>It seems you've gone somewhere you shouldn't!  404 NOT FOUND!</p>
	<p/>
	<p>This can happen if you enter the spotify id wrong and search, go back and check it!</p>
		<p class="lead">Spotify IDs look like this : 37i9dQZF1DX4UtSsGT1Sbe - 22 characters and can be got by clicking share and "Copy (Playlist|Album) Link".</p>
		<p>Shameless self-promotion : Follow me on twitch - <a href="https://www.twitch.tv/iMartynOnTwitch">iMartynOnTwitch</a>, oddly enough, I do a lot of twitchsings!</p>
  </main>
`)
	sendAllTheRest(response)
}

func CSSHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-type","text/css")
	fmt.Fprint(response,
`
/*
* Globals
*/

/* Links */
a,
a:focus,
a:hover {
 color: #fff;
}

/* Custom default button */
.btn-secondary,
.btn-secondary:hover,
.btn-secondary:focus {
 color: #333;
 text-shadow: none; /* Prevent inheritance from body */
 background-color: #fff;
 border: .05rem solid #fff;
}


/*
* Base structure
*/

html,
body {
 height: 100%;
 background-color: #333;
}

body {
 display: -ms-flexbox;
 display: flex;
 color: #fff;
 text-shadow: 0 .05rem .1rem rgba(0, 0, 0, .5);
 box-shadow: inset 0 0 5rem rgba(0, 0, 0, .5);
}

.cover-container {
 max-width: 142em;
}


/*
* Header
*/
.masthead {
 margin-bottom: 2rem;
}

.masthead-brand {
 margin-bottom: 0;
}

.nav-masthead .nav-link {
 padding: .25rem 0;
 font-weight: 700;
 color: rgba(255, 255, 255, .5);
 background-color: transparent;
 border-bottom: .25rem solid transparent;
}

.nav-masthead .nav-link:hover,
.nav-masthead .nav-link:focus {
 border-bottom-color: rgba(255, 255, 255, .25);
}

.nav-masthead .nav-link + .nav-link {
 margin-left: 1rem;
}

.nav-masthead .active {
 color: #fff;
 border-bottom-color: #fff;
}

@media (min-width: 48em) {
 .masthead-brand {
   float: left;
 }
 .nav-masthead {
   float: right;
 }
}


/*
* Cover
*/
.cover {
 padding: 0 1.5rem;
}
.cover .btn-lg {
 padding: .75rem 1.25rem;
 font-weight: 700;
}


/*
* Footer
*/
.mastfoot {
 color: rgba(255, 255, 255, .5);
}
`)
}

func PlaylistHandler(response http.ResponseWriter, request *http.Request) {
	type listItem struct {
		TrackName     string
		SpotifyArtist string
		TwitchArtist  string
		MatchKind     MatchType
	}
	listItems := []listItem{}
	vars := mux.Vars(request)
	if len(cachedSongList) <= 0 {
		fmt.Println("No twitch sings cached yet, reading from disk or remote...")
		CachedTwitchGetSongs(true)
	} else {
		CachedTwitchGetSongs(false)
	}
	playlistTracks, playlistInfo := SpotifyGetPlaylistTracks(vars["playlist"])
	foundCount := 0
	for _,spotifySong := range(playlistTracks) {
		item := listItem{}
		item.TrackName = spotifySong.Name
		item.SpotifyArtist = SpotifyArtistsAsString(spotifySong.Artists)
		artists := []string{}
		for _, artistforname := range spotifySong.Artists {
			artists = append(artists, artistforname.Name)
		}
		ret,twitchSong := SpotifyListContains(item.TrackName,artists)
		item.MatchKind = ret
		if ret != MatchNoMatch {
			item.TwitchArtist = twitchSong.Artist
			foundCount += 1
		}
		listItems = append(listItems, item)
	}
	sendHTMLPreambleAndHead(response)
	sendBodyAndHeader(response)
	fmt.Fprint(response,
`<main role="main" class="inner cover">
	<h1 class="cover-heading">Songs in your playlist "`+playlistInfo.Name+`": </h1>
	<div align="left" style="width: 74%; display: inline-block; padding-left: 40px;">`+strconv.Itoa(playlistInfo.TrackCount)+` tracks, `+strconv.Itoa(foundCount)+` (possibly) in Twitch sings catalog of `+strconv.Itoa(len(cachedSongList))+`!</div>
	<div align="right" style="width: 24%; display: inline-block;"><a href="#" id="showhidebutton" onclick="javascript:toggleUnfound()">Toggle unmatched songs</a></div>
	<ul>
`)
	for _,item := range(listItems) {
		matchTypeString := ""
		twitchArtistString := ""
		matchClass := ""
		switch item.MatchKind {
			case MatchNoMatch : 
				matchTypeString = "[NO MATCH]"
				matchClass = "match-nomatch"
			case MatchTrackNameOnly : 
				matchTypeString = "[MATCH SONG NAME]"
				matchClass = "match-matchtrack"
				twitchArtistString = " ("+item.TwitchArtist+")"
			case MatchBothNameAndArtist : 
				matchTypeString = "[MATCH SONG AND ARTIST]"
				matchClass = "match-fullmatch"
			case MatchNameOnlyFuzzy : 
				matchTypeString = "[MATCH SONG NAME FUZZY]"
				matchClass = "match-matchtrackfuzzt"
				twitchArtistString = " ("+item.TwitchArtist+")"
			case MatchBothFuzzy : 
				matchTypeString = "[MATCH SONG AND ARTIST FUZZY]"
				matchClass = "match-fullmatchfuzzy"
		}
		if item.MatchKind != MatchNoMatch {
		}
		fmt.Fprintf(response,"<li class=\""+matchClass+"\">%s - %s%s %s</li>",item.TrackName,item.SpotifyArtist,twitchArtistString,matchTypeString)
	}
	fmt.Fprint(response, `</ul></main>`)
	sendAllTheRest(response)
}

func AlbumHandler(response http.ResponseWriter, request *http.Request) {
	type listItem struct {
		TrackName     string
		SpotifyArtist string
		TwitchArtist  string
		MatchKind     MatchType
	}
	listItems := []listItem{}
	vars := mux.Vars(request)
	if len(cachedSongList) <= 0 {
		fmt.Println("No twitch sings cached yet, reading from disk or remote...")
		CachedTwitchGetSongs(true)
	} else {
		CachedTwitchGetSongs(false)
	}
	albumTracks, albumInfo := SpotifyGetAlbumTracks(vars["album"])
	foundCount := 0
	for _,spotifySong := range(albumTracks) {
		item := listItem{}
		item.TrackName = spotifySong.Name
		item.SpotifyArtist = SpotifyArtistsAsString(spotifySong.Artists)
		artists := []string{}
		for _, artistforname := range spotifySong.Artists {
			artists = append(artists, artistforname.Name)
		}
		ret,twitchSong := SpotifyListContains(item.TrackName,artists)
		item.MatchKind = ret
		if ret != MatchNoMatch {
			item.TwitchArtist = twitchSong.Artist
			foundCount += 1
		}
		listItems = append(listItems, item)
	}
	sendHTMLPreambleAndHead(response)
	sendBodyAndHeader(response)
	fmt.Fprint(response,
`<main role="main" class="inner cover">
	<h1 class="cover-heading">Songs in your album "`+albumInfo.Name+`": </h1>
	<div align="left" style="width: 74%; display: inline-block; padding-left: 40px;">`+strconv.Itoa(len(albumTracks))+` tracks, `+strconv.Itoa(foundCount)+` (possibly) in Twitch sings catalog of `+strconv.Itoa(len(cachedSongList))+`!</div>
	<div align="right" style="width: 24%; display: inline-block;"><a href="#" id="showhidebutton" onclick="javascript:toggleUnfound()">Toggle unmatched songs</a></div>
	<ul>
`)
	for _,item := range(listItems) {
		matchTypeString := ""
		twitchArtistString := ""
		matchClass := ""
		switch item.MatchKind {
			case MatchNoMatch : 
				matchTypeString = "[NO MATCH]"
				matchClass = "match-nomatch"
			case MatchTrackNameOnly : 
				matchTypeString = "[MATCH SONG NAME]"
				matchClass = "match-matchtrack"
				twitchArtistString = " ("+item.TwitchArtist+")"
			case MatchBothNameAndArtist : 
				matchTypeString = "[MATCH SONG AND ARTIST]"
				matchClass = "match-fullmatch"
			case MatchNameOnlyFuzzy : 
				matchTypeString = "[MATCH SONG NAME FUZZY]"
				matchClass = "match-matchtrackfuzzt"
				twitchArtistString = " ("+item.TwitchArtist+")"
			case MatchBothFuzzy : 
				matchTypeString = "[MATCH SONG AND ARTIST FUZZY]"
				matchClass = "match-fullmatchfuzzy"
		}
		if item.MatchKind != MatchNoMatch {
		}
		fmt.Fprintf(response,"<li class=\""+matchClass+"\">%s - %s%s %s</li>",item.TrackName,item.SpotifyArtist,twitchArtistString,matchTypeString)
	}
	fmt.Fprint(response, `</ul></main>`)
	sendAllTheRest(response)
}

func HandleHTTP() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/cover.css", CSSHandler)
	r.HandleFunc(`/playlist/{playlist:[0-9a-zA-Z]{22}}`, PlaylistHandler)
	r.HandleFunc(`/playlist/spotify:playlist:{playlist:[0-9a-zA-Z]{22}}`, PlaylistHandler)
	r.HandleFunc("/album/{album:[0-9a-zA-Z]{22}}", AlbumHandler)
	r.HandleFunc("/album/spotify:album:{album:[0-9a-zA-Z]{22}}", AlbumHandler)
	http.Handle("/", r)
	srv := &http.Server {
		Handler: r,
		Addr: "0.0.0.0:5353",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
	}
	srv.ListenAndServe()
}

