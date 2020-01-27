package dao

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"gochatroom/chat"
	"log"
)

// hash:  	gochat.room    key  --> json string
// hash:    gochat.member  key  --> json string
const (
	// 聊天室hash key
	roomKey = "gochat:room"
	// 用户hash key
	memberKey = "gochat:member"
	// 聊天室-成员set前缀
	roomMemberSetKeyPrefix = "gochat:room-member:"
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
	client.HSet(roomKey, roomId, string(marshal))
}

// 获取指定聊天室信息
func RetrieveRoom(roomId string) chat.Room {
	var room chat.Room
	result := client.Get(roomId).String()
	err := json.Unmarshal([]byte(result), &room)
	if err != nil {
		log.Println("反序列化room结果时出错", err)
	}
	return room
}

// 获取所有聊天室信息
func RetrieveAllRooms() []chat.Room {
	var rooms []chat.Room
	result, err := client.HGetAll(roomKey).Result()
	if err != nil {
		log.Println("获取room信息出错", err)
	}
	for _, str := range result {
		var room chat.Room
		err := json.Unmarshal([]byte(str), &room)
		if err != nil {
			log.Println("反序列化room结果时出错", err)
		}
		rooms = append(rooms, room)
	}
	return rooms
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
	var member chat.Member
	result := client.HGet(memberKey, memberId).String()
	err := json.Unmarshal([]byte(result), &member)
	if err != nil {
		log.Println("反序列化member结果时出错", err)
	}
	return member
}

// 获取所有聊天室信息
func RetrieveAllMembers() []chat.Member {
	var members []chat.Member
	result, err := client.HGetAll(memberKey).Result()
	if err != nil {
		log.Println("获取member信息出错", err)
	}
	for _, str := range result {
		var member chat.Member
		err := json.Unmarshal([]byte(str), &member)
		if err != nil {
			log.Println("反序列化member结果时出错", err)
		}
		members = append(members, member)
	}
	return members
}

// 聊天成员进入聊天室
func MemberRoomEnter(roomId, memberId string)  {
	key := roomMemberSetKeyPrefix + roomId
	client.SAdd(key, memberId)
}

// 聊天成员退出聊天室
func MemberRoomExit(roomId, memberId string)  {
	key := roomMemberSetKeyPrefix + roomId
	client.SRem(key, memberId)
}

// 获取聊天室所有成员id
func RetrieveRoomAllMembers(roomId string) []string {
	key := roomMemberSetKeyPrefix + roomId
	result, err := client.SMembers(key).Result()
	if err != nil {
		log.Println("获取聊天室成员信息出错", err)
	}
	return result
}