package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type MessageType int

const (
	MemberMsg MessageType = iota
	MemberJoin
	MemberExit
)

const LISTEN_ADDR  = ":8080"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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

type Member struct {
	id       string
	nickname string
	conn     *websocket.Conn
	msgCh    chan Message
	endCh	 chan int
	wg		 sync.WaitGroup
	room     *ChatRoom
}

func (m *Member) notifyExit()  {
	m.endCh <- 0
}

func (m *Member) ListenMessage() {
	defer func() {
		m.room.exit <- m
		if err := m.conn.Close(); err != nil {
			log.Println("关闭ws连接失败", err)
		}
	}()

	m.wg.Add(2)
	// 启动两个goroutine持续监听写和读
	go m.keepWriting()
	go m.keepReading()
	// 等待两个goroutine都退出后再继续执行
	m.wg.Wait()
}

func (m *Member) keepReading()  {
	defer m.wg.Done()

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

func (m *Member) keepWriting()  {
	defer m.wg.Done()

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

type ChatRoom struct {
	members     map[*Member]bool
	join        chan *Member
	exit        chan *Member
	bcast       chan Message
	memberCount int
}

func (r *ChatRoom) broadcast(msg Message) {
	r.bcast <- msg
}

func (r *ChatRoom) Run() {
	for {
		select {
		case mem := <-r.join:
			r.members[mem] = true
			go r.broadcast(NewMessage(MemberJoin, mem.nickname, ""))
		case mem := <-r.exit:
			if _, ok := r.members[mem]; ok {
				delete(r.members, mem)
				close(mem.msgCh)
				go r.broadcast(NewMessage(MemberExit, mem.nickname, ""))
			}
		case msg := <-r.bcast:
			for mem := range r.members {
				select {
				case mem.msgCh <- msg:
				default:
					delete(r.members, mem)
					close(mem.msgCh)
				}
			}

		}
	}
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		members:     make(map[*Member]bool),
		join:        make(chan *Member),
		exit:        make(chan *Member),
		bcast:       make(chan Message),
		memberCount: 0,
	}
}

func HandleWsRequest(cr *ChatRoom, w http.ResponseWriter, r *http.Request)  {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("获取ws连接出错", err)
		return
	}

	name := r.FormValue("name")
	mem := &Member{
		id:       uuid.New().String(),
		nickname: name,
		conn:     conn,
		msgCh:    make(chan Message, 128),
		endCh:	  make(chan int),
		room:     cr,
	}

	mem.room.join <- mem

	go mem.ListenMessage()
}

func main() {
	log.Println("服务器启动")
	room := NewChatRoom()
	go room.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWsRequest(room, w, r)
	})

	if err := http.ListenAndServe(LISTEN_ADDR, nil); err != nil {
		log.Fatalln("主程序出错，程序退出", err)
	}
}
