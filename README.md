 # spotifytwitchsings
 
 A tool to take a spotify playlist and compare it to the twitchsings songs to see if there's any you can sing!

----

This is now running at https://spotifytwitch.martyn.berlin/

## Installation
Make sure you have golang installed.
```bash
go get github.com/imartyn/spotifytwitchsings
```


## Usage

```bash
# help

spotifytwitchsings -p 13rgfJS9aI8PwfuDCaGJp0
# Voila!!
```

or

```bash
spotifytwitchsings serve
# then it is listening on 0.0.0.0:5353
```

## Notes

* This is very rough around the edges, think of it like a hackday project (yes I did it in ~12 hrs).
* Windows build is pretty certain not to work
* There's so much that I should do but this is a little side-project that got out of hand because I wanted friends to be able to use it!
    * Port is hardcoded
    * Cache is hardcoded
    * Not a lot of input validation
    * Some pretty important refactoring needed
    * ...but hey, it works, so ship it!

## Credits

Shout out to [spotifydl](https://github.com/BharatKalluri/spotifydl) by BharatKalluri which I hacked about to make this. 