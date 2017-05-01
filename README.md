# YTSbot

Telegram bot to automate the search and download of movies (from [yts.ag](https://yts.ag)) using [Transmission](https://transmissionbt.com/).

## Dependencies
- Transmission web cli
## Run

Talk with the Botfather to create a new bot. Set the following commands:<br>
random_yts - Shows an inline keyboard to select a movie category. If a number is specified, sets the page.<br>
search_yts - Looks for a specified movie.<br>
bt_list - Lists Torrents.<br>
bt_del - Removes the specified Torrent.<br>

Then
```
$ go get github.com/0-wHiTeHand-0/ytsbot
$ vim $GOPATH/src/github.com/0-wHiTeHand-0/ytsbot/bot.cfg //Set the config
$ $GOPATH/bin/ytsbot $GOPATH/src/github.com/0-wHiTeHand-0/ytsbot/bot.cfg
```
