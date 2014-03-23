# confreaks

[![Build Status](https://drone.io/github.com/subosito/confreaks/status.png)](https://drone.io/github.com/subosito/confreaks/latest)

confreaks on the command line


## Installation

Grab binary from [release page](https://github.com/subosito/confreaks/releases/latest)

Or if you have go installed, you can simply:

```
go get github.com/subosito/confreaks
```

## Usage

Create a dedicated directory for saving confreaks' conferences.

```
$ mkdir conferences
$ cd conferences
```

Run `confreaks init` to get all events from confreaks.com. It will saves to an `index.json` file.

```
$ confreaks init
```

Now you can print all the events using `confreaks list`:

```
$ confreaks list
In Memory of Jim Weirich
Big Ruby 2014
DevconTLV January 2014
.....
Ruby Conference 2007
Ruby Hoedown 2007
MountainWest RubyConf 2007
```

Once you have list of events, now you can syncronize per event or all events. You can use `confreaks sync` command.

```
$ # synchronize all events, beware this takes moments to finished
$ confreaks sync

$ # synchronize per event
$ confreaks sync "ArrrrCamp 2013"

$ # synchronize all events that contains "2013"
$ confreaks sync "2013"
```

Finally, you can download the presentation videos. Just like above `sync` command, you can use `confreaks download`.

```
$ confreaks download "ArrrrCamp 2013"
2014/03/23 13:57:18 ++ ArrrrCamp 2013
2014/03/23 13:57:18  +-- Downloading You gotta try this
[youtube] Setting language
[youtube] sVd4p6oKeUA: Downloading webpage
[youtube] sVd4p6oKeUA: Downloading video info webpage
[youtube] sVd4p6oKeUA: Extracting video information
[download] Destination: ArrrrCamp 2013/ArrrrCamp 2013 - You gotta try this by Avdi Grimm-sVd4p6oKeUA.mp4
[download]   0.4% of 106.19MiB at 304.79KiB/s ETA 05:55
```

## Dependencies

- [youtube-dl](https://github.com/rg3/youtube-dl): For downloading presentation videos.

## Tips

By default `youtube-dl` will download best quality videos. If you want to get lower quality or specific quality, you can set on `~/.config/youtube-dl.conf` with, eg:

```
-f '18/43'
```

Which means `youtube-dl` will only download video with format:

```
18 mp4  640x360
43 webm 640x360
```

You can find another format by using `youtube-dl -F VIDEO_URL`.

