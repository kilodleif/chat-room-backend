package chat

import "time"

type MessageType int

const (
	MemberMsg MessageType = iota
	MemberJoin
	MemberExit
)

type Message struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	Content   string      `json:"content"`
	Timestamp string      `json:"timestamp"`
}

func NewMessage(msgType MessageType, nickname, content string) Message {
	return Message{
		Type:      msgType,
		From:      nickname,
		Content:   content,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}
}
