package chat

import (
	"github.com/google/uuid"
	"time"
)

const maxMemberCountPerRoom = 20

type Room struct {
	id          string
	maxJoins    int
	curJoins    int
	createTime  string
	lastActTime string
}

type RoomControl struct {
	room        *Room
	join        chan *MemberControl
	exit        chan *MemberControl
	broadcast   chan Message
	memControls map[*MemberControl]bool
}

func NewRoomControl() *RoomControl {
	nowStr := time.Now().Format(TimeFormat)
	return &RoomControl{
		room: &Room{
			id:          uuid.New().String(),
			maxJoins:    maxMemberCountPerRoom,
			curJoins:    0,
			createTime:  nowStr,
			lastActTime: nowStr,
		},
		memControls: make(map[*MemberControl]bool),
		join:        make(chan *MemberControl),
		exit:        make(chan *MemberControl),
		broadcast:   make(chan Message),
	}
}

func (c *RoomControl) Run() {
	go func(r *RoomControl) {
		for {
			select {
			case mCtrl := <-r.join:
				r.memControls[mCtrl] = true
				go r.broadcastMessage(
					NewMessage(
						MemberJoin,
						mCtrl.member.nickname,
						"",
					),
				)
			case mCtrl := <-r.exit:
				if _, ok := r.memControls[mCtrl]; ok {
					delete(r.memControls, mCtrl)
					close(mCtrl.msgCh)
					go r.broadcastMessage(
						NewMessage(
							MemberExit,
							mCtrl.member.nickname,
							"",
						),
					)
				}
			case msg := <-r.broadcast:
				for mCtrl := range r.memControls {
					select {
					case mCtrl.msgCh <- msg:
					default:
						delete(r.memControls, mCtrl)
						close(mCtrl.msgCh)
					}
				}

			}
		}
	}(c)
}

func (c *RoomControl) broadcastMessage(msg Message) {
	c.broadcast <- msg
}
