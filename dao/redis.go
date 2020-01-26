package dao

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"gochatroom/chat"
	"log"
	"time"
)

// string:  gochat.room.<room_id>      --> json string
// hash:    gochat.member  key  --> json string
const (
	// 聊天室string key前缀
	roomKeyPrefix = "gochat:room:"
	// 用户hash key
	memberKey = "gochat:member"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// 保存聊天室信息
func SaveRoom(roomId string, room *chat.Room)  {
	marshal, err := json.Marshal(room)
	if err != nil {
		log.Println("序列化room时出错", err)
	}
	key := roomKeyPrefix + roomId
	client.Set(key, string(marshal), 36 * time.Hour)
}

// 获取指定聊天室信息
func RetrieveRoom(roomId string) chat.Room {
	result := client.Get(roomId).String()
	room := &chat.Room{}
	err := json.Unmarshal([]byte(result), room)
	if err != nil {
		log.Println("反序列化room结果时出错", err)
	}
	return *room
}

// 保存用户信息
func SaveMember(memberId string, member *chat.Member)  {
	marshal, err := json.Marshal(member)
	if err != nil {
		log.Println("序列化member时出错", err)
	}
	client.HSet(memberKey, memberId, string(marshal))
}

// 获取指定用户信息
func RetrieveMember(memberId string) chat.Member {
	result := client.HGet(memberKey, memberId).String()
	member := &chat.Member{}
	err := json.Unmarshal([]byte(result), member)
	if err != nil {
		log.Println("反序列化member结果时出错", err)
	}
	return *member
}
