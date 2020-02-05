package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

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

type JoinRoomReq struct {
	ID       string `json:"id"`
	FromID   string `json:"fromId"`
	RoomID   string `json:"roomId"`
	StreamID string `json:"streamId"`
	EndType  string `json:"endType"`
}

type JoinRoomReply struct {
	ID             string   `json:"id"`
	MasterID       string   `json:"masterId"`
	OnlineUserList []string `json:"onlineUserList"`
}

func (r JoinRoomReply) ToBytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

type PushMessageReq struct {
	ConversationType int    `json:"conversationType"`
	TargetId         string `json:"targetId"`
	FromId           string `json:"fromId"`
	Content          string `json:"content"`
	ObjectName       string `json:"objectName"`
}

func (p PushMessageReq) ToPushToClient() PushToClient {
	return PushToClient{
		ConversationType: p.ConversationType,
		TargetId:         p.TargetId,
		FromId:           p.FromId,
		Content:          p.Content,
		ObjectName:       p.ObjectName,
		SentTime:         time.Now().Unix(),
		MessageUId:       uuid.New().String(),
	}
}

type PushRoomReq struct {
	Room string `json:"room"`
	Msg  string `json:"msg"`
}

type PushToClient struct {
	ConversationType int    `json:"conversationType"`
	TargetId         string `json:"targetId"`
	FromId           string `json:"fromId"`
	Content          string `json:"content"`
	ObjectName       string `json:"objectName"`
	SentTime         int64  `json:"sentTime"`
	MessageUId       string `json:"messageUId"`
}

func (p PushToClient) ToBytes() []byte {
	b, _ := json.Marshal(p)
	return b
}
