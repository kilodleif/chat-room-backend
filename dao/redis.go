package dao

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"gochatroom/chat"
	"log"
)

// gochat.room.<room_id>      --> json string
// gochat.member.<member_id>  --> json string
const (
	// 聊天室redis key前缀
	roomKeyPrefix = "gochat.room."
	// 用户redis key前缀
	memberKeyPrefix = "gochat.member."
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
	client.RPush(key, string(marshal))
}

// 获取所有聊天室信息
func RetrieveRooms(roomId string) []chat.Room {
	var rooms []chat.Room
	key := roomKeyPrefix + roomId
	result, err := client.LRange(key, 0, -1).Result()
	if err != nil {
		log.Println("获取room结果时出错", err)
	}
	for _, str := range result {
		room := &chat.Room{}
		err := json.Unmarshal([]byte(str), room)
		if err != nil {
			log.Println("反序列化room结果时出错", err)
		}
		rooms = append(rooms, *room)
	}
	return rooms
}

// 保存用户信息
func SaveMember(memberId string, member *chat.Member)  {
	marshal, err := json.Marshal(member)
	if err != nil {
		log.Println("序列化member时出错", err)
	}
	key := memberKeyPrefix + memberId
	client.RPush(key, string(marshal))
}

// 获取所有用户信息
func RetrieveMembers(memberId string) []chat.Member {
	var members []chat.Member
	key := memberKeyPrefix + memberId
	result, err := client.LRange(key, 0, -1).Result()
	if err != nil {
		log.Println("获取member结果时出错", err)
	}
	for _, str := range result {
		member := &chat.Member{}
		err := json.Unmarshal([]byte(str), member)
		if err != nil {
			log.Println("反序列化member结果时出错", err)
		}
		members = append(members, *member)
	}
	return members
}
