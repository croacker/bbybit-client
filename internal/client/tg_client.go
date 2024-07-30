package client

import (
	"log"

	"github.com/croacker/bybit-client/internal/config"
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

func StartBot() {
	bot, err := tg_bot_api.NewBotAPI("") // TODO bot token
	if err != nil {
		panic(err)
	}

	go readIncoming(bot)
	go writeOutgoing(bot)
}

func readIncoming(bot *tg_bot_api.BotAPI) {
	upd := tg_bot_api.NewUpdate(0)
	upd.Timeout = 30

	updates := bot.GetUpdatesChan(upd)

	for update := range updates {
		id := getChatId(update)
		if id != -1 {
			saveChatId(id)
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
		log.Println("out msg:", msg)
		for id, _ := range tgClient.chatIds {
			log.Println("send msg to:", id)
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

func saveChatId(id int64) {
	if _, ok := tgClient.chatIds[id]; !ok {
		tgClient.chatIds[id] = id
	}
}
