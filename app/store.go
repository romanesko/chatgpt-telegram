package main

import (
	"encoding/json"
	"git.mills.io/prologic/bitcask"
)

type ChatsList struct {
	Chats map[int64]void `json:"chats"`
}

func getChatsListStore() ChatsList {
	db, _ := bitcask.Open("/tmp/data/db")
	defer db.Close()
	val, _ := db.Get([]byte("allowedChats"))
	var values map[int64]void
	json.Unmarshal(val, &values)
	return ChatsList{values}
}

func (c *ChatsList) addChat(chatId int64) {
	if c.Chats == nil {
		c.Chats = map[int64]void{}
	}
	c.Chats[chatId] = void{}
	db, _ := bitcask.Open("/tmp/data/db")
	defer db.Close()
	encodedValues, _ := json.Marshal(c.Chats)
	db.Put([]byte("allowedChats"), encodedValues)
}

func (c *ChatsList) isChatAllowed(chatId int64) bool {
	_, exists := c.Chats[chatId]
	return exists
}

func (c *ChatsList) listChats() []int64 {
	var keys []int64
	for k := range c.Chats {
		keys = append(keys, k)
	}
	return keys
}
