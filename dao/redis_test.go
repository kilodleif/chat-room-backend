package dao

import (
	"github.com/google/uuid"
	"gochatroom/chat"
	"reflect"
	"testing"
)

var (
	room1 chat.Room
	member1 chat.Member
	timestampString = "2020-01-26 11:22:33"
)

func TestMain(m *testing.M) {
	// 初始化 测试用例
	room1 = chat.Room{
		Id:          uuid.New().String(),
		MaxJoins:    10,
		CurJoins:    0,
		CreateTime:  timestampString,
		LastActTime: timestampString,
	}

	member1 = chat.Member{
		Id:       uuid.New().String(),
		Nickname: "Bob",
		JoinTime: timestampString,
	}

	m.Run()
	// 测试完成后 删除redis中的数据
	client.Del(roomKey, memberKey)
}

func TestSaveRoom(t *testing.T) {
	type args struct {
		roomId string
		room   *chat.Room
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"保存Room用例1",
			args{
				roomId: room1.Id,
				room:   &room1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveRoom(tt.args.roomId, tt.args.room)
		})
	}
}

func TestRetrieveRoom(t *testing.T) {
	type args struct {
		roomId string
	}
	tests := []struct {
		name string
		args args
		want chat.Room
	}{
		// TODO: Add test cases.
		{
			"获取Room用例1",
			args{roomId:room1.Id},
			room1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetrieveRoom(tt.args.roomId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSaveMember(t *testing.T) {
	type args struct {
		memberId string
		member   *chat.Member
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"保存Member用例1",
			args{
				memberId: member1.Id,
				member:   &member1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SaveMember(tt.args.memberId, tt.args.member)
		})
	}
}

func TestRetrieveMember(t *testing.T) {
	type args struct {
		memberId string
	}
	tests := []struct {
		name string
		args args
		want chat.Member
	}{
		// TODO: Add test cases.
		{
			"获取Member用例1",
			args{memberId:member1.Id},
			member1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetrieveMember(tt.args.memberId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RetrieveMember() = %v, want %v", got, tt.want)
			}
		})
	}
}