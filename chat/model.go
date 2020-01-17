package chat

import "time"

type MessageType int

const (
	MemberMsg MessageType = iota
	MemberJoin
	MemberExit
)

const TimeFormat  = "2006-01-02 15:04:05"

type Message struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	Content   string      `json:"content"`
	Timestamp string      `json:"timestamp"`
}

func NewMessage(msgType MessageType, from, content string) Message {
	return Message{
		Type:      msgType,
		From:      from,
		Content:   content,
		Timestamp: time.Now().Format(TimeFormat),
	}
}
