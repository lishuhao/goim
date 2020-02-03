package model

import "encoding/json"

const (
	UserId1 = "haiping"
	UserId2 = "yuyao"
	UserId3 = "wenshuai"
	UserId4 = "shuhao"
)
const (
	Mid1 = 1001
	Mid2 = 1002
	Mid3 = 1003
	Mid4 = 1004
)

var Uid2Mid = map[string]int64{
	UserId1: Mid1,
	UserId2: Mid2,
	UserId3: Mid3,
	UserId4: Mid4,
}

type CreateRoomReq struct {
	ID       string `json:"id"`
	FromID   string `json:"fromId"`
	RoomID   string `json:"roomId"`
	StreamID string `json:"streamId"`
	EndType  string `json:"endType"`
}

type CreateRoomReply struct {
	ID             string   `json:"id"`
	MasterID       string   `json:"masterId"`
	OnlineUserList []string `json:"onlineUserList"`
}

func (r CreateRoomReply) ToBytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

type PushKeysReq struct {
	Keys []string `json:"keys"`
	Msg  string   `json:"msg"`
}

type PushRoomReq struct {
	Room string `json:"room"`
	Msg  string `json:"msg"`
}

type PushAllReq struct {
	Msg string `json:"msg"`
}
