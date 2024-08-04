package db

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/croacker/bybit-client/internal/config"
	"go.etcd.io/bbolt"
)

// var db *bbolt.
var dbPath string

func SetupDb(cfg *config.AppConfig) {
	dbPath = cfg.DbCfg.Path
	log.Println("setup db:", dbPath)
	db := OpenDb()
	defer db.Close()

	err := db.Update(func(tx *bbolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket:%v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("TG_CHAT"))
		if err != nil {
			return fmt.Errorf("could not create tg_chat bucket:%v", err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Println("db setup done")
}

func OpenDb() *bbolt.DB {
	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func SaveChat(chat TgChat) {
	db := OpenDb()
	defer db.Close()

	k := []byte(strconv.FormatInt(chat.Id, 10))
	v, err := json.Marshal(chat)
	if err != nil {
		log.Fatal("error marshal chat:", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		err = tx.Bucket([]byte("DB")).Put(k, v)
		if err != nil {
			return fmt.Errorf("error save chat:%v", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func AllChats() {
	db := OpenDb()
	defer db.Close()

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("DB")).Bucket([]byte("TG_CHAT"))
		b.ForEach(func(k, v []byte) error {
			fmt.Println(string(k), string(v))
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
