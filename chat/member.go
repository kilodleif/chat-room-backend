package chat

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const messageChannelBuffer = 128

type Member struct {
	id       string
	nickname string
	joinTime string
}

type MemberControl struct {
	conn     *websocket.Conn
	msgCh    chan Message
	endCh    chan int
	member   *Member
	roomCtrl *RoomControl
}

func NewMemberControl(nickname string, conn *websocket.Conn, roomCtrl *RoomControl) *MemberControl {
	return &MemberControl{
		conn:  conn,
		msgCh: make(chan Message, messageChannelBuffer),
		endCh: make(chan int),
		member: &Member{
			id:       uuid.New().String(),
			nickname: nickname,
			joinTime: time.Now().Format(TimeFormat),
		},
		roomCtrl: roomCtrl,
	}
}

func (c *MemberControl) ListenMessage() {
	go func(c *MemberControl) {
		defer func(c *MemberControl) {
			c.roomCtrl.exit <- c
			if err := c.conn.Close(); err != nil {
				log.Println("关闭ws连接失败", err)
			}
		}(c)
		// 通知新成员加入事件
		c.roomCtrl.join <- c

		var wg sync.WaitGroup
		wg.Add(2)
		// 启动两个goroutine持续监听写和读
		go c.keepWriting(&wg)
		go c.keepReading(&wg)
		// 等待两个goroutine都退出后再继续执行
		wg.Wait()
	}(c)
}

func (c *MemberControl) keepReading(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-c.endCh:
			return
		default:
			resp := struct {
				Message string `json:"msg"`
			}{}
			if err := c.conn.ReadJSON(&resp); err != nil {
				log.Println("读取消息失败，goroutine退出", err)
				// 退出前通知 keepWriting goroutine
				c.notifyExit()
				return
			}
			c.roomCtrl.broadcastMessage(NewMessage(MemberMsg, c.member.nickname, resp.Message))
		}
	}
}

func (c *MemberControl) keepWriting(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-c.endCh:
			return
		case msg := <-c.msgCh:
			if err := c.conn.WriteJSON(msg); err != nil {
				log.Println("写入消息失败，goroutine退出", err)
				// 退出前通知 keepReading goroutine
				c.notifyExit()
				return
			}
		}
	}
}

func (c *MemberControl) notifyExit() {
	c.endCh <- 0
}
