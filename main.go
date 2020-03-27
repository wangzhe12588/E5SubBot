package main

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/http"
	"time"
)

var (
	BotToken string
	Socks5   string
	bot      *tb.Bot
	db       *sql.DB
)

const (
	dbDriverName = "sqlite3"
	dbName       = "./data.db"
)

func init() {
	//read config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	CheckErr(err)
	BotToken = viper.GetString("bot_token")
	Socks5 = viper.GetString("socks5")
	//set bot
	botsettings := tb.Settings{
		Token:  BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}
	//set socks5
	if Socks5 != "" {
		fmt.Println("Proxy:" + Socks5)
		dialer, err := proxy.SOCKS5("tcp", Socks5, nil, proxy.Direct)
		CheckErr(err)
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		httpTransport.Dial = dialer.Dial
		botsettings.Client = httpClient
	}
	//create bot
	bot, err = tb.NewBot(botsettings)
	CheckErr(err)

	//sqlite init
	db, err = sql.Open(dbDriverName, dbName)
	CheckErr(err)
	CreateTB(db)
}
func main() {
	BotStart()
	//b.Handle(tb.OnText, func(m *tb.Message) {
	//	b.Send(m.Sender, "hello world")
	//})

	//b.Start()
}
func BotStart() {
	makeHandle()
	bot.Start()
}
func makeHandle() {
	bot.Handle("/start", bStart)
	bot.Handle("/my", bMy)
	bot.Handle("/bind", bBind)
	bot.Handle("/about", bAbout)
	bot.Handle(tb.OnText, bOnText)
	//bot.Handle(tb.InlineButton{Unique: ""})
}
