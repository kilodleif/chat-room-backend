package chat

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const messageChannelBuffer = 128

type Member struct {
	id       string
	nickname string
	conn     *websocket.Conn
	msgCh    chan Message
	endCh    chan int
	room     *Room
}

func NewMember(nickname string, conn *websocket.Conn, room *Room) *Member {
	return &Member{
		id:       uuid.New().String(),
		nickname: nickname,
		conn:     conn,
		msgCh:    make(chan Message, messageChannelBuffer),
		endCh:    make(chan int),
		room:     room,
	}
}

func (m *Member) ListenMessage() {
	go func(m *Member) {
		defer func(m *Member) {
			m.room.exit <- m
			if err := m.conn.Close(); err != nil {
				log.Println("关闭ws连接失败", err)
			}
		}(m)
		// 通知新成员加入事件
		m.room.join <- m

		var wg sync.WaitGroup
		wg.Add(2)
		// 启动两个goroutine持续监听写和读
		go m.keepWriting(&wg)
		go m.keepReading(&wg)
		// 等待两个goroutine都退出后再继续执行
		wg.Wait()
	}(m)
}

func (m *Member) keepReading(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-m.endCh:
			return
		default:
			resp := struct {
				Message string `json:"msg"`
			}{}
			if err := m.conn.ReadJSON(&resp); err != nil {
				log.Println("读取消息失败，goroutine退出", err)
				// 退出前通知 keepWriting goroutine
				m.notifyExit()
				return
			}
			m.room.broadcast(NewMessage(MemberMsg, m.nickname, resp.Message))
		}
	}
}

func (m *Member) keepWriting(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-m.endCh:
			return
		case msg := <-m.msgCh:
			if err := m.conn.WriteJSON(msg); err != nil {
				log.Println("写入消息失败，goroutine退出", err)
				// 退出前通知 keepReading goroutine
				m.notifyExit()
				return
			}
		}
	}
}

func (m *Member) notifyExit() {
	m.endCh <- 0
}
