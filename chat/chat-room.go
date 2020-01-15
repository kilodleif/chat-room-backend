package chat

type Room struct {
	members     map[*Member]bool
	join        chan *Member
	exit        chan *Member
	bcast       chan Message
	memberCount int
}

func NewRoom() *Room {
	return &Room{
		members:     make(map[*Member]bool),
		join:        make(chan *Member),
		exit:        make(chan *Member),
		bcast:       make(chan Message),
		memberCount: 0,
	}
}

func (r *Room) Run() {
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

func (r *Room) broadcast(msg Message) {
	r.bcast <- msg
}
