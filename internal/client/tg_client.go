package client

import (
	"log"

	"github.com/croacker/bybit-client/internal/config"
	"github.com/croacker/bybit-client/internal/db"
	tg_bot_api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var tgClient *TgBotClient

type TgBotClient struct {
	token        string
	outgoingChan chan string
	chatIds      map[int64]int64
}

func NewTgClient(cfg *config.AppConfig) *TgBotClient {
	tgClient = &TgBotClient{
		cfg.TgClient.Token,
		make(chan string),
		make(map[int64]int64),
	}
	return tgClient
}

func (t *TgBotClient) GetOutgoingChannel() chan string {
	return t.outgoingChan
}

func (t *TgBotClient) Start() {
	log.Println("start tg-client...")
	bot, err := tg_bot_api.NewBotAPI(t.token)
	if err != nil {
		panic(err)
	}

	go readIncoming(bot)
	go writeOutgoing(bot)
	log.Println("tg-client started")
}

func readIncoming(bot *tg_bot_api.BotAPI) {
	upd := tg_bot_api.NewUpdate(0)
	upd.Timeout = 30

	updates := bot.GetUpdatesChan(upd)

	for update := range updates {
		id := getChatId(update)
		if id != -1 {
			saveChatId(update)

			if update.Message != nil {
				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				msg := tg_bot_api.NewMessage(id, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
		}
	}
}

func writeOutgoing(bot *tg_bot_api.BotAPI) {
	for msg := range tgClient.outgoingChan {
		for _, chat := range db.AllChats() {
			id := chat.Id
			tgMsg := tg_bot_api.NewMessage(id, msg)
			_, err := bot.Send(tgMsg)
			if err != nil {
				log.Println("error send message to:", chat)
			}
		}
	}
}

func getChatId(update tg_bot_api.Update) int64 {
	var result int64 = -1
	if update.Message != nil {
		result = update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		result = update.CallbackQuery.Message.Chat.ID
	}

	return result
}

func getChatInfo(update tg_bot_api.Update) *tg_bot_api.Chat {
	var result *tg_bot_api.Chat
	if update.Message != nil {
		result = update.Message.Chat
	}

	if update.CallbackQuery != nil {
		result = update.CallbackQuery.Message.Chat
	}

	return result
}

func saveChatId(update tg_bot_api.Update) {
	chat := getChatInfo(update)
	dbChat := db.TgChat{
		chat.ID,
		chat.Type,
		chat.UserName,
		chat.FirstName,
		chat.LastName,
	}
	db.SaveChat(dbChat)
}
