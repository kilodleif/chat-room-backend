package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"gochatroom/chat"
	"log"
	"net/http"
)

const DefaultListenAddr = ":8080"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var addr = flag.String("addr", DefaultListenAddr, "address to listen")

func main() {
	flag.Parse()
	log.Println("服务器启动，监听", *addr)
	room := chat.NewRoomControl()
	room.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWsRequest(room, w, r)
	})

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln("主程序出错，程序退出", err)
	}
}

func handleWsRequest(room *chat.RoomControl, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("获取ws连接出错", err)
		return
	}

	name := r.FormValue("name")
	mem := chat.NewMemberControl(name, conn, room)
	mem.ListenMessage()
}
