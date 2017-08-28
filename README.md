# YTSbot

Telegram bot to automate the search and download of movies (from [yts.ag](https://yts.ag)) using [Transmission](https://transmissionbt.com/).

## Dependencies
- Transmission web cli
## Run

Talk with the Botfather to create a new bot. Set the following commands:<br><br>
random_yts - Shows an inline keyboard to select a movie category. If a number is specified, sets the page.<br>
search_yts - Looks for a specified movie.<br>
bt_list - Lists torrents.<br>
bt_clean - Removes completed torrents from the list (not the data).<br>
bt_del - Removes the specified torrent (data included).<br>

Then
```
$ go get github.com/0-wHiTeHand-0/ytsbot
$ vim $GOPATH/src/github.com/0-wHiTeHand-0/ytsbot/bot.cfg //Set the config
$ $GOPATH/bin/ytsbot $GOPATH/src/github.com/0-wHiTeHand-0/ytsbot/bot.cfg
```
Don't forget to include the IP of the host where runs the bot, in the Transmission RPC whitelist.
