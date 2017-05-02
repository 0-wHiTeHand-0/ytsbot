package main

import (
      "regexp"
      "strconv"
      "strings"
      "errors"
      "github.com/lnguyen/go-transmission/transmission"
      )
const (
      StatusStopped = iota
      StatusCheckPending
      StatusChecking
      StatusDownloadPending
      StatusDownloading
      StatusSeedPending
      StatusSeeding
      )

type CmdConfigTransmission struct {
   St_list        string
      Reg_del      *regexp.Regexp
}

var bt_client  transmission.TransmissionClient

func NewCmdTransmission(config *CmdConfigTransmission, url string, user string, pass string) {
   config.St_list = "/bt_list"
      config.Reg_del = regexp.MustCompile(`^/bt_del [0-9]+$`)
      bt_client = transmission.New(url, user, pass)
}


func BtDown (hash string) (string){
   torrent, err := bt_client.AddTorrentByURL("https://yts.ag/torrent/download/" + hash, "")
      var stemp string
      if err != nil{
         stemp = "Error!\n\n" + err.Error()
      }else if torrent.Name == ""{
         stemp = "Error!\n\nYou don't have permission to download this torrent. Try manually, but you will probably need to be logged in."
         }else{
         stemp = "Success!\n\nAdded: " + torrent.Name + "\nID: " + strconv.Itoa(torrent.ID)
      }
   return stemp
}

func BtList() (string, error){
   torrents, err := bt_client.GetTorrents()
      if err != nil{
         return "", err
      }
      if len(torrents) == 0{
         return "No torrents in the list!", nil
      }
      var s_eta string
s := "<--- Torrent list --->\n"
      for _, torrent := range torrents{
         if torrent.Eta == -1{
            s_eta = "Not available"
            }else if torrent.Eta == -2{
               s_eta = "Unknown"
            }else{
               s_eta = strconv.Itoa(torrent.Eta/60) + " Min"
            }
         s += "\n<" + strconv.Itoa(torrent.ID) + "> " + torrent.Name + "\n" + BtStatus(&torrent) +
         " - Remains " + strconv.Itoa(torrent.LeftUntilDone/1000000) + "MB - Complete " +
         strconv.FormatFloat(torrent.PercentDone*100, 'f', -1, 32) + "% - ETA " +
         s_eta
      }
      return s, nil
}

func BtStatus (t *transmission.Torrent) string {
   switch t.Status {
      case StatusStopped:
         return "Stopped"
      case StatusCheckPending:
            return "Check waiting"
      case StatusChecking:
               return "Checking"
      case StatusDownloadPending:
                  return "Download waiting"
      case StatusDownloading:
                     return "Downloading"
      case StatusSeedPending:
                        return "Seed waiting"
      case StatusSeeding:
                           return "Seeding"
      default:
                              return "unknown"
   }

}

func BtRemove(a string) (string, error){
   b := strings.SplitN(a, " ", 2)
   if len(b) != 2{
      return "", errors.New("Error parsing bt_del command")
   }
   c, err := strconv.Atoi(b[1])
   if err != nil{
      return "", err
   }
   d, err := bt_client.RemoveTorrent(c, true)//Removes the torrent AND the file
   return d, err
}
