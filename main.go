package main

import (
      "os"
      "io/ioutil"
      "encoding/json"
      "log"
      "github.com/go-telegram-bot-api/telegram-bot-api"
      )

type d_transmission struct{
   Url   string   `json:"url"`
   Username string   `json:"username"`
   Password string   `json:"password"`
}

type config struct {
   Token          string     `json:"token"`
   AllowedIDs     []int      `json:"allowed_ids"`
   Transmission   d_transmission     `json:"transmission"`
}

type cmdConfigs struct {
   Yts   CmdConfigYts
   Transmission   CmdConfigTransmission
}

var Commandssl cmdConfigs

func main() {
   if len(os.Args) != 2 {
      log.Fatalln("Usage: ./ytsbot <CONFIGFILE>")
   }
   cfg, err := parseConfig(os.Args[1])
   if err != nil{
      log.Fatalln(err)
   }
   bot, err := tgbotapi.NewBotAPI(cfg.Token)
      if err != nil {
         log.Panic(err)
      }
      NewCmdYts(&Commandssl.Yts)
      NewCmdTransmission(&Commandssl.Transmission, cfg.Transmission.Url,cfg.Transmission.Username,cfg.Transmission.Password)

      bot.Debug = false
      log.Printf("Authorized on account %s", bot.Self.UserName)

      u := tgbotapi.NewUpdate(0)
      u.Timeout = 60

      updates, err := bot.GetUpdatesChan(u)
      for update := range updates {
         go handle_updates(bot, update, cfg)
      }
}

func handle_updates(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg config) {
   if update.Message!=nil{
      if !allowed_ids(update.Message.From.ID, cfg.AllowedIDs){
        log.Println("User " + update.Message.From.FirstName + " blocked! Alias: " + update.Message.From.UserName)
        return
        }
      log.Printf("Message -> [%s] %s (id: %d, alias: %s)", update.Message.From.FirstName, update.Message.Text, update.Message.Chat.ID, update.Message.From.UserName)
         var msg tgbotapi.MessageConfig
         if Commandssl.Yts.Reg_rand.MatchString(update.Message.Text){
            YtsSetPage(update.Message.Text)
               msg = YtsNewKeyboard(update.Message.Chat.ID)
         }else if Commandssl.Yts.Reg_detail.MatchString(update.Message.Text){
            img, m, torrent, err := YtsDetail(update.Message.Text)
               if err != nil{
                  log.Println(err)
                     return
               }
            if img != ""{
               msg = tgbotapi.NewMessage(update.Message.Chat.ID, img)
                  msg.DisableWebPagePreview = false
                  bot.Send(msg)
         }
         msg = tgbotapi.NewMessage(update.Message.Chat.ID, m)
            if torrent != ""{
               var DownloadKeyboard = tgbotapi.NewInlineKeyboardMarkup(
                     tgbotapi.NewInlineKeyboardRow(
                        tgbotapi.NewInlineKeyboardButtonData("Download", "download_torrent_" + torrent),
                        ),
                     )
                  msg.ReplyMarkup = DownloadKeyboard
         }
         }else if Commandssl.Yts.Reg_search.MatchString(update.Message.Text){
            m, err := YtsSearch(update.Message.Text)
               if err != nil{
                  log.Println(err)
                     return
               }
            msg = tgbotapi.NewMessage(update.Message.Chat.ID, m)
         }else if Commandssl.Transmission.St_list == update.Message.Text{
            m, err := BtList()
               if err != nil{
                  log.Println(err)
                     return
               }
            msg = tgbotapi.NewMessage(update.Message.Chat.ID, m)
         }else if Commandssl.Transmission.Reg_del.MatchString(update.Message.Text){
            m, err := BtRemove(update.Message.Text)
               if err != nil{
                  log.Println(err)
                     return
               }
            msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Operation result: " + m)
         }else{
m := "Command not found!"
      msg = tgbotapi.NewMessage(update.Message.Chat.ID, m)
      msg.ReplyToMessageID = update.Message.MessageID
         }
      msg.DisableWebPagePreview = false
         bot.Send(msg)
   }else if update.CallbackQuery != nil {
      if !allowed_ids(update.CallbackQuery.From.ID, cfg.AllowedIDs){
        log.Println("User " + update.CallbackQuery.From.FirstName + " blocked! Alias: " + update.CallbackQuery.From.UserName)
        return
        }
callback := tgbotapi.CallbackConfig{
CallbackQueryID: update.CallbackQuery.ID,
          }
          log.Printf("Callback -> [%s] %s (id: %d, alias: %s)", update.CallbackQuery.From.FirstName,
                update.CallbackQuery.Data, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.From.UserName)
             if ((len(update.CallbackQuery.Data) > 3) && (update.CallbackQuery.Data[:4] == "YTS_")){
                m, err := YtsRandom(update.CallbackQuery.Data[4:])
                   if err != nil{
                      log.Println(err)
                         return
                   }
msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, m)
        var MoreKeyboard = tgbotapi.NewInlineKeyboardMarkup(
              tgbotapi.NewInlineKeyboardRow(
                 tgbotapi.NewInlineKeyboardButtonData("More...", update.CallbackQuery.Data),
                 ),
              )
        msg.ReplyMarkup = MoreKeyboard
        bot.Send(msg)
             }else if ((len(update.CallbackQuery.Data) > 20) && (update.CallbackQuery.Data[:17] == "download_torrent_")){
                /*var DoneKeyboard = tgbotapi.NewInlineKeyboardMarkup(
                  tgbotapi.NewInlineKeyboardRow(
                  tgbotapi.NewInlineKeyboardButtonData("Done!", ""),
                  ),
                  )
                //msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID,update.CallbackQuery.Message.MessageID , DoneKeyboard)
msg :=
tgbotapi.NewEditMessageCaption(update.CallbackQuery.Message.Chat.ID,update.CallbackQuery.Message.MessageID,"LEL")
bot.Send(msg)*/
                callback.Text = "Downloading movie..." + BtDown(update.CallbackQuery.Data[17:])
                   callback.ShowAlert = true
             }
          if _, err := bot.AnswerCallbackQuery(callback); err != nil {
             log.Println(err)
          }
   }else{
      log.Println("Unknown query received: ", update)
   }
}

func YtsNewKeyboard(id int64) tgbotapi.MessageConfig {

   var CategoryKeyboard = tgbotapi.NewInlineKeyboardMarkup(
         tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Action", "YTS_Action"),
            tgbotapi.NewInlineKeyboardButtonData("Adventure", "YTS_Adventure"),
            tgbotapi.NewInlineKeyboardButtonData("Animation", "YTS_Animation"),
            ),
         tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Comedy", "YTS_Comedy"),
            tgbotapi.NewInlineKeyboardButtonData("Documentary", "YTS_Documentary"),
            tgbotapi.NewInlineKeyboardButtonData("Drama", "YTS_Drama"),
            ),
         tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Fantasy", "YTS_Fantasy"),
            tgbotapi.NewInlineKeyboardButtonData("Horror", "YTS_Horror"),
            tgbotapi.NewInlineKeyboardButtonData("Mystery", "YTS_Mystery"),
            ),
         tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("Sport", "YTS_Sport"),
            tgbotapi.NewInlineKeyboardButtonData("Sci-Fi", "YTS_Sci-fi"),
            tgbotapi.NewInlineKeyboardButtonData("Thriller", "YTS_Thriller"),
            ),
         tgbotapi.NewInlineKeyboardRow(
               tgbotapi.NewInlineKeyboardButtonData("All", "YTS_"),
               ),
         )
            msg := tgbotapi.NewMessage(id, "Select a category:")
            msg.ReplyMarkup = CategoryKeyboard
            return msg
}

func parseConfig(file string) (config, error) {
   b, err := ioutil.ReadFile(file)
      if err != nil {
         return config{}, err
      }
   var cfg config
      if err := json.Unmarshal(b, &cfg); err != nil {
         return config{}, err
      }
   log.Println("Config: ", cfg)
      return cfg, nil
}

func allowed_ids(a int, b []int) bool{
   for _, i := range b{
      if i == a{return true}
   }
   return false
}
