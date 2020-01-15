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
	room     *ChatRoom
}

func (m *Member) ListenMessage() {
	var wg sync.WaitGroup
	wg.Add(2)

	defer func() {
		m.room.exit <- m
		if err := m.conn.Close(); err != nil {
			log.Println("EXIT PROC", err)
		}
	}()

	go m.keepWriting(wg)
	go m.keepReading(wg)

	wg.Wait()
}

func (m *Member) keepReading(wg sync.WaitGroup)  {
	defer wg.Done()

	for {
		resp := struct {
			Message string `json:"msg"`
		}{}
		if err := m.conn.ReadJSON(&resp); err != nil {
			log.Println("READ PROC", err)
			return
		}
		m.room.broadcast(NewMessage(MemberMsg, m.nickname, resp.Message))
	}
}

func (m *Member) keepWriting(wg sync.WaitGroup)  {
	defer wg.Done()

	for {
		msg := <-m.msgCh
		if err := m.conn.WriteJSON(msg); err != nil {
			log.Println("WRITE PROC", err)
			return
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
			// ???
			go r.broadcast(NewMessage(MemberJoin, mem.nickname,  mem.nickname + " has joined"))
		case mem := <-r.exit:
			if _, ok := r.members[mem]; ok {
				delete(r.members, mem)
				close(mem.msgCh)
				// ???
				go r.broadcast(NewMessage(MemberExit, mem.nickname, mem.nickname + " has left"))
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
		return
	}

	name := r.FormValue("name")
	mem := &Member{
		id:       uuid.New().String(),
		nickname: name,
		conn:     conn,
		msgCh:    make(chan Message, 128),
		room:     cr,
	}
	mem.room.join <- mem

	go mem.ListenMessage()
}

func main() {
	log.Println("server started")
	room := NewChatRoom()
	go room.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		HandleWsRequest(room, w, r)
	})

	if err := http.ListenAndServe(LISTEN_ADDR, nil); err != nil {
		log.Fatalln("MAIN PROC", err)
	}
}
